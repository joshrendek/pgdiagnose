package pgdiagnose

import (
	"bytes"
	"github.com/kr/logfmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func getData(logplexURL string) []byte {
	resp, err := http.Get(logplexURL)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil
	}

	return body
}

func findLogsForDatabase(logs []byte, database string) []byte {
	lines := bytes.Split(logs, []byte("\n"))
	search := []byte("source=" + database)
	for i := len(lines) - 1; i >= 0; i-- {
		if bytes.Contains(lines[i], search) {
			return lines[i]
		}
	}
	return nil
}

type DatabaseLog struct {
	LoadAvg1M      float64 `logfmt:"sample#load-avg-1m"`
	LoadAvg5M      float64 `logfmt:"sample#load-avg-5m"`
	LoadAvg15M     float64 `logfmt:"sample#load-avg-15m"`
	MemoryTotal    int64
	MemoryCached   int64
	MemoryFree     int64
	MemoryPostgres int64

	RawMemoryTotal    string `logfmt:"sample#memory-total"`
	RawMemoryCached   string `logfmt:"sample#memory-cached"`
	RawMemoryFree     string `logfmt:"sample#memory-free"`
	RawMemoryPostgres string `logfmt:"sample#memory-postgres"`
}

func parseLog(log []byte) (parsed DatabaseLog) {
	logfmt.Unmarshal(log, &parsed)
	parsed.MemoryTotal = removeKb(parsed.RawMemoryTotal)
	parsed.MemoryCached = removeKb(parsed.RawMemoryTotal)
	parsed.MemoryFree = removeKb(parsed.RawMemoryTotal)
	parsed.MemoryPostgres = removeKb(parsed.RawMemoryPostgres)
	return parsed
}

func removeKb(raw string) int64 {
	foo, _ := strconv.ParseInt(strings.TrimRight(raw, "kB"), 10, 0)
	return foo
}
