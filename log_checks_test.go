package main

import (
	"bytes"
	"testing"
)

var (
	ivory = findLogForDatabase("HEROKU_POSTGRESQL_IVORY", sampleData)
	red   = findLogForDatabase("HEROKU_POSTGRESQL_RED", sampleData)
)

func TestCheckLoadOnLog(t *testing.T) {
	checkIvory := checkLoadOnLog(ivory)
	if checkIvory.Status != "green" {
		t.Fatalf("Ivory wasn't green")
	}

	checkRed := checkLoadOnLog(red)
	if checkRed.Status != "red" {
		t.Fatalf("Red wasn't red")
	}

	checkMissing := checkLoadOnLog(DatabaseLog{})
	if checkMissing.Status != "skipped" {
		t.Fatalf("Missing wasn't skipped")
	}

}

func TestFindLogLineForDatabase(t *testing.T) {
	line := findLogLineForDatabase("HEROKU_POSTGRESQL_IVORY", sampleData)
	if !bytes.HasPrefix(line, []byte("2014-05-22T20:30:21+00:00")) {
		t.Fatalf("got wrong line: %v", string(line))
	}

	line = findLogLineForDatabase("LOL", sampleData)
	if !(line == nil) {
		t.Fatalf("should not have found anything, got line: %v", string(line))
	}
}

func TestParseLog(t *testing.T) {
	parsed := parseLog(findLogLineForDatabase("HEROKU_POSTGRESQL_IVORY", sampleData))
	expected := DatabaseLog{"HEROKU_POSTGRESQL_IVORY", 0.17, 0.16, 0.215,
		7629452, 6930188, 165524, 205648,
		"7629452kB", "6930188kB", "165524kB", "205648kB"}
	if parsed != expected {
		t.Fatalf("got %v instead of %v", parsed, expected)
	}

	parsed = parseLog([]byte(""))
	expected = DatabaseLog{}
	if parsed != expected {
		t.Fatalf("got %v instead of %v", parsed, expected)
	}
}

var sampleData = []byte(`014-05-22T20:29:48+00:00 app[heroku-postgres]: source=HEROKU_POSTGRESQL_IVORY sample#current_transaction=1879 sample#db_size=6801592bytes sample#tables=1 sample#active-connections=2 sample#waiting-connections=0 sample#index-cache-hit-rate=0.71429 sample#table-cache-hit-rate=0.75 sample#load-avg-1m=0.07 sample#load-avg-5m=0.16 sample#load-avg-15m=0.215 sample#read-iops=17.058 sample#write-iops=22.251 sample#memory-total=7629452kB sample#memory-free=165524kB sample#memory-cached=6930188kB sample#memory-postgres=405648kB
2014-05-22T20:30:21+00:00 app[heroku-postgres]: source=HEROKU_POSTGRESQL_IVORY sample#current_transaction=1879 sample#db_size=6801592bytes sample#tables=1 sample#active-connections=2 sample#waiting-connections=0 sample#index-cache-hit-rate=0.71429 sample#table-cache-hit-rate=0.75 sample#load-avg-1m=0.17 sample#load-avg-5m=0.16 sample#load-avg-15m=0.215 sample#read-iops=17.058 sample#write-iops=22.251 sample#memory-total=7629452kB sample#memory-free=165524kB sample#memory-cached=6930188kB sample#memory-postgres=205648kB
2014-05-22T20:30:53+00:00 app[heroku-postgres]: source=HEROKU_POSTGRESQL_RED sample#current_transaction=1879 sample#db_size=6801592bytes sample#tables=1 sample#active-connections=2 sample#waiting-connections=0 sample#index-cache-hit-rate=0.71429 sample#table-cache-hit-rate=0.75 sample#load-avg-1m=6.405 sample#load-avg-5m=0.17 sample#load-avg-15m=0.215 sample#read-iops=34.117 sample#write-iops=21.931 sample#memory-total=7629452kB sample#memory-free=132448kB sample#memory-cached=6963204kB sample#memory-postgres=205648kB`)
