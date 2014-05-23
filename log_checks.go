package pgdiagnose

import (
	"bytes"
	"github.com/kr/logfmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func CheckLogs(database string, logplexURL string) []Check {
	return checkLogsFromBytes(database, getData(logplexURL))
}

func checkLogsFromBytes(database string, data []byte) []Check {
	log := findLogForDatabase(database, data)

	v := make([]Check, 6)
	v[0] = checkLoadOnLog(log)

	return v
}

type loadAvgs struct {
	LoadAvg1M  float64
	LoadAvg5M  float64
	LoadAvg15M float64
}

func checkLoadOnLog(log DatabaseLog) Check {
	if (log == DatabaseLog{}) {
		return Check{"Load", "skipped", nil}
	} else {
		load := loadAvgs{log.LoadAvg1M, log.LoadAvg5M, log.LoadAvg15M}
		return checkLoad(load)
	}
}

func checkLoad(load loadAvgs) Check {
	if load.LoadAvg1M > 3 {
		return Check{"Load", "red", load}
	} else if load.LoadAvg1M > 1 {
		return Check{"Load", "yellow", load}
	} else {
		return Check{"Load", "green", nil}
	}
}

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

func findLogLineForDatabase(database string, logs []byte) []byte {
	lines := bytes.Split(logs, []byte("\n"))
	search := []byte("source=" + database)
	for i := len(lines) - 1; i >= 0; i-- {
		if bytes.Contains(lines[i], search) {
			return lines[i]
		}
	}
	return nil
}

func findLogForDatabase(database string, logs []byte) DatabaseLog {
	return parseLog(findLogLineForDatabase(database, logs))
}

type DatabaseLog struct {
	Source         string  `logfmt:"source="`
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
	parsed.MemoryCached = removeKb(parsed.RawMemoryCached)
	parsed.MemoryFree = removeKb(parsed.RawMemoryFree)
	parsed.MemoryPostgres = removeKb(parsed.RawMemoryPostgres)
	return parsed
}

func removeKb(raw string) int64 {
	foo, _ := strconv.ParseInt(strings.TrimRight(raw, "kB"), 10, 0)
	return foo
}
