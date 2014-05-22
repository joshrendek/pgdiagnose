package pgdiagnose

import (
	"bytes"
	"testing"
)

func TestFindLogsForDatabase(t *testing.T) {
	line := findLogsForDatabase(sampleData, "HEROKU_POSTGRESQL_IVORY")
	if !bytes.HasPrefix(line, []byte("2014-05-22T20:30:21+00:00")) {
		t.Fatalf("got wrong line: %v", string(line))
	}
	line = findLogsForDatabase(sampleData, "LOL")
	if !(line == nil) {
		t.Fatalf("should not have found anything, got line: %v", string(line))
	}
}

func TestParseLog(t *testing.T) {
	parsed := parseLog(sampleData)
	expected := DatabaseLog{0.205, 0.17, 0.215, 7629452, 7629452, 7629452, 205648, "7629452kB", "6963204kB", "132448kB", "205648kB"}
	if parsed != expected {
		t.Fatalf("got %v instead of %v", parsed, expected)
	}
}

var sampleData = []byte(`014-05-22T20:29:48+00:00 app[heroku-postgres]: source=HEROKU_POSTGRESQL_IVORY sample#current_transaction=1879 sample#db_size=6801592bytes sample#tables=1 sample#active-connections=2 sample#waiting-connections=0 sample#index-cache-hit-rate=0.71429 sample#table-cache-hit-rate=0.75 sample#load-avg-1m=0.07 sample#load-avg-5m=0.16 sample#load-avg-15m=0.215 sample#read-iops=17.058 sample#write-iops=22.251 sample#memory-total=7629452kB sample#memory-free=165524kB sample#memory-cached=6930188kB sample#memory-postgres=205648kB
2014-05-22T20:30:21+00:00 app[heroku-postgres]: source=HEROKU_POSTGRESQL_IVORY sample#current_transaction=1879 sample#db_size=6801592bytes sample#tables=1 sample#active-connections=2 sample#waiting-connections=0 sample#index-cache-hit-rate=0.71429 sample#table-cache-hit-rate=0.75 sample#load-avg-1m=0.07 sample#load-avg-5m=0.16 sample#load-avg-15m=0.215 sample#read-iops=17.058 sample#write-iops=22.251 sample#memory-total=7629452kB sample#memory-free=165524kB sample#memory-cached=6930188kB sample#memory-postgres=205648kB
2014-05-22T20:30:53+00:00 app[heroku-postgres]: source=HEROKU_POSTGRESQL_RED sample#current_transaction=1879 sample#db_size=6801592bytes sample#tables=1 sample#active-connections=2 sample#waiting-connections=0 sample#index-cache-hit-rate=0.71429 sample#table-cache-hit-rate=0.75 sample#load-avg-1m=0.205 sample#load-avg-5m=0.17 sample#load-avg-15m=0.215 sample#read-iops=34.117 sample#write-iops=21.931 sample#memory-total=7629452kB sample#memory-free=132448kB sample#memory-cached=6963204kB sample#memory-postgres=205648kB`)
