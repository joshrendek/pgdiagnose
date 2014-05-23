package main

import (
	"testing"
)

var sanitizetests = []struct {
	in  JobParams
	out JobParams
}{
	{JobParams{}, JobParams{}},
	{JobParams{"", "", "", "", ""}, JobParams{"", "", "", "", ""}},
	{JobParams{"postgres://user:pass@host.com/db", "", "", "", ""}, JobParams{"postgres://user:pass@host.com/db", "", "", "", ""}},
	{JobParams{"", "", "crane", "", ""}, JobParams{"", "", "crane", "", ""}},
	{JobParams{"", "", "cr@ne", "", ""}, JobParams{"", "", "", "", ""}},
	{JobParams{"", "", "", "sushi", ""}, JobParams{"", "", "", "sushi", ""}},
	{JobParams{"", "", "", "su$hi", ""}, JobParams{"", "", "", "", ""}},
	{JobParams{"", "", "", "", "HEROKU_POSTGRESQL_RED_URL"}, JobParams{"", "", "", "", "HEROKU_POSTGRESQL_RED_URL"}},
	{JobParams{"", "", "", "", "&EROKU_POSTGRESQL_RED_URL"}, JobParams{"", "", "", "", ""}},
}

func TestSanitizeJopParams(t *testing.T) {
	for i, tt := range sanitizetests {
		tt.in.sanitize()
		if tt.in != tt.out {
			t.Errorf("%d. Expected to sanitize to %v, but was %v", i, tt.out, tt.in)
		}
	}
}
