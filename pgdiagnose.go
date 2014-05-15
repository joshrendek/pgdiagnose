package main

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type check struct {
	Name    string
	Results interface{}
}

func main() {
	connstring := "dbname=will sslmode=disable"
	if len(os.Args) > 1 {
		connstring = os.Args[1]
	}
	db := connectDB(connstring)

	v := make([]check, 5)

	v[0] = check{"Long Queries", longQueriesCheck(db)}
	v[1] = check{"Idle in Transaction", idleQueriesCheck(db)}
	v[2] = check{"Unused Indexes", unusedIndexesCheck(db)}
	v[3] = check{"Bloat", bloatCheck(db)}
	v[4] = check{"Hit Rate", hitRateCheck(db)}
	js, _ := json.Marshal(v)
	fmt.Println("what: ", string(js))
}

func errDie(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func connectDB(dbURL string) *sqlx.DB {
	db, err := sqlx.Open("postgres", dbURL)
	errDie(err)

	_, err = db.Exec("select 1")
	errDie(err)

	return db
}

type longQueriesResult struct {
	Pid      int64
	Duration float64
	Query    string
}

func longQueriesCheck(db *sqlx.DB) (results []longQueriesResult) {
	query := `
	  SELECT pid, now()-query_start as duration, query
	  FROM pg_stat_activity
	  WHERE now()-query_start > '1 minute'::interval
		AND state <> 'idle in transaction'
		;`
	err := db.Select(&results, query)
	errDie(err)
	return results
}

func longQueriesStatus(results []longQueriesResult) string {
	if len(results) == 0 {
		return "green"
	} else {
		return "red"
	}
}

type idleQueriesResult struct {
	Pid      int64
	Duration float64
	Query    string
}

func idleQueriesCheck(db *sqlx.DB) (results []idleQueriesResult) {
	query := `
	  SELECT pid, now()-query_start as duration, query
	  FROM pg_stat_activity
	  WHERE now()-query_start > '1 minute'::interval
		AND state like 'idle in trans%'
		;`
	err := db.Select(&results, query)
	errDie(err)
	return results
}

func idleQueriesStatus(results []idleQueriesResult) string {
	if len(results) == 0 {
		return "green"
	} else {
		return "red"
	}
}

type unusedIndexesResult struct {
	Reason          string
	Schemaname      string
	Tablename       string
	Indexname       string
	Index_scan_pct  string
	Scans_per_write string
	Index_size      string
	Table_size      string
}

func unusedIndexesCheck(db *sqlx.DB) (results []unusedIndexesResult) {
	// http://www.databasesoup.com/2014/05/new-finding-unused-indexes-query.html
	query := `
WITH table_scans as (
    SELECT relid,
        tables.idx_scan + tables.seq_scan as all_scans,
        ( tables.n_tup_ins + tables.n_tup_upd + tables.n_tup_del ) as writes,
                pg_relation_size(relid) as table_size
        FROM pg_stat_user_tables as tables
),
all_writes as (
    SELECT sum(writes) as total_writes
    FROM table_scans
),
indexes as (
    SELECT idx_stat.relid, idx_stat.indexrelid,
        idx_stat.schemaname, idx_stat.relname as tablename,
        idx_stat.indexrelname as indexname,
        idx_stat.idx_scan,
        pg_relation_size(idx_stat.indexrelid) as index_bytes,
        indexdef ~* 'USING btree' AS idx_is_btree
    FROM pg_stat_user_indexes as idx_stat
        JOIN pg_index
            USING (indexrelid)
        JOIN pg_indexes as indexes
            ON idx_stat.schemaname = indexes.schemaname
                AND idx_stat.relname = indexes.tablename
                AND idx_stat.indexrelname = indexes.indexname
    WHERE pg_index.indisunique = FALSE
),
index_ratios AS (
SELECT schemaname, tablename, indexname,
    idx_scan, all_scans,
    round(( CASE WHEN all_scans = 0 THEN 0.0::NUMERIC
        ELSE idx_scan::NUMERIC/all_scans * 100 END),2) as index_scan_pct,
    writes,
    round((CASE WHEN writes = 0 THEN idx_scan::NUMERIC ELSE idx_scan::NUMERIC/writes END),2)
        as scans_per_write,
    pg_size_pretty(index_bytes) as index_size,
    pg_size_pretty(table_size) as table_size,
    idx_is_btree, index_bytes
    FROM indexes
    JOIN table_scans
    USING (relid)
),
index_groups AS (
SELECT 'Never Used Indexes' as reason, *, 1 as grp
FROM index_ratios
WHERE
    idx_scan = 0
    and idx_is_btree
UNION ALL
SELECT 'Low Scans, High Writes' as reason, *, 2 as grp
FROM index_ratios
WHERE
    scans_per_write <= 1
    and index_scan_pct < 10
    and idx_scan > 0
    and writes > 100
    and idx_is_btree
UNION ALL
SELECT 'Seldom Used Large Indexes' as reason, *, 3 as grp
FROM index_ratios
WHERE
    index_scan_pct < 5
    and scans_per_write > 1
    and idx_scan > 0
    and idx_is_btree
    and index_bytes > 100000000
UNION ALL
SELECT 'High-Write Large Non-Btree' as reason, index_ratios.*, 4 as grp
FROM index_ratios, all_writes
WHERE
    ( writes::NUMERIC / ( total_writes + 1 ) ) > 0.02
    AND NOT idx_is_btree
    AND index_bytes > 100000000
ORDER BY grp, index_bytes DESC )
SELECT reason, schemaname, tablename, indexname,
    index_scan_pct, scans_per_write, index_size, table_size
FROM index_groups;
`
	err := db.Select(&results, query)
	errDie(err)
	return results
}

func unusedIndexesStatus(results []unusedIndexesResult) string {
	if len(results) == 0 {
		return "green"
	} else {
		return "red"
	}
}

type bloatResult struct {
	Type   string
	Object string
	Bloat  float64
	Waste  string
}

func bloatCheck(db *sqlx.DB) (results []bloatResult) {
	query := `
WITH constants AS (
  SELECT current_setting('block_size')::numeric AS bs, 23 AS hdr, 4 AS ma
), bloat_info AS (
  SELECT
    ma,bs,schemaname,tablename,
    (datawidth+(hdr+ma-(case when hdr%ma=0 THEN ma ELSE hdr%ma END)))::numeric AS datahdr,
    (maxfracsum*(nullhdr+ma-(case when nullhdr%ma=0 THEN ma ELSE nullhdr%ma END))) AS nullhdr2
  FROM (
    SELECT
      schemaname, tablename, hdr, ma, bs,
      SUM((1-null_frac)*avg_width) AS datawidth,
      MAX(null_frac) AS maxfracsum,
      hdr+(
        SELECT 1+count(*)/8
        FROM pg_stats s2
        WHERE null_frac<>0 AND s2.schemaname = s.schemaname AND s2.tablename = s.tablename
      ) AS nullhdr
    FROM pg_stats s, constants
    GROUP BY 1,2,3,4,5
  ) AS foo
), table_bloat AS (
  SELECT
    schemaname, tablename, cc.relpages, bs,
    CEIL((cc.reltuples*((datahdr+ma-
      (CASE WHEN datahdr%ma=0 THEN ma ELSE datahdr%ma END))+nullhdr2+4))/(bs-20::float)) AS otta
  FROM bloat_info
  JOIN pg_class cc ON cc.relname = bloat_info.tablename
  JOIN pg_namespace nn ON cc.relnamespace = nn.oid AND nn.nspname = bloat_info.schemaname AND nn.nspname <> 'information_schema'
), index_bloat AS (
  SELECT
    schemaname, tablename, bs,
    COALESCE(c2.relname,'?') AS iname, COALESCE(c2.reltuples,0) AS ituples, COALESCE(c2.relpages,0) AS ipages,
    COALESCE(CEIL((c2.reltuples*(datahdr-12))/(bs-20::float)),0) AS iotta
  FROM bloat_info
  JOIN pg_class cc ON cc.relname = bloat_info.tablename
  JOIN pg_namespace nn ON cc.relnamespace = nn.oid AND nn.nspname = bloat_info.schemaname AND nn.nspname <> 'information_schema'
  JOIN pg_index i ON indrelid = cc.oid
  JOIN pg_class c2 ON c2.oid = i.indexrelid
)
SELECT
  type, object, bloat, pg_size_pretty(raw_waste) as waste
FROM
(SELECT
  'table' as type,
  schemaname ||'.'|| tablename as object,
  ROUND(CASE WHEN otta=0 THEN 0.0 ELSE table_bloat.relpages/otta::numeric END,1) AS bloat,
  CASE WHEN relpages < otta THEN '0' ELSE (bs*(table_bloat.relpages-otta)::bigint)::bigint END AS raw_waste
FROM
  table_bloat
    UNION
SELECT
  'index' as type,
  schemaname || '.' || tablename || '::' || iname as object,
  ROUND(CASE WHEN iotta=0 OR ipages=0 THEN 0.0 ELSE ipages/iotta::numeric END,1) AS bloat,
  CASE WHEN ipages < iotta THEN '0' ELSE (bs*(ipages-iotta))::bigint END AS raw_waste
FROM
  index_bloat) bloat_summary
WHERE raw_waste > 10*1024*1024 AND bloat > 10
ORDER BY raw_waste DESC, bloat DESC
;`
	err := db.Select(&results, query)
	errDie(err)
	return results
}

func bloatStatus(results []bloatResult) string {
	if len(results) == 0 {
		return "green"
	} else {
		return "red"
	}
}

type hitRateResult struct {
	Name  string
	Ratio float64
}

func hitRateCheck(db *sqlx.DB) (results []hitRateResult) {
	query := `
WITH rates AS (
	SELECT
		'index hit rate' AS name,
		sum(idx_blks_hit) / nullif(sum(idx_blks_hit + idx_blks_read), 0) AS ratio
	FROM pg_statio_user_indexes
	UNION ALL
	SELECT
		'table hit rate' AS name,
		sum(heap_blks_hit) / nullif(sum(heap_blks_hit) + sum(heap_blks_read), 0) AS ratio
	FROM pg_statio_user_tables
)
SELECT * FROM rates WHERE ratio < 0.99
;`
	err := db.Select(&results, query)
	errDie(err)
	return results
}

func hitRateStatus(results []hitRateResult) string {
	if len(results) == 0 {
		return "green"
	} else {
		return "red"
	}
}
