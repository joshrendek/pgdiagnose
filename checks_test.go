package pgdiagnose

import (
	"testing"
)

func TestLongQueriesStatus(t *testing.T) {
	values := make([]longQueriesResult, 0)
	if longQueriesStatus(values) != "green" {
		t.Fatal("not green on empty results")
	}

	values = make([]longQueriesResult, 1)
	if longQueriesStatus(values) != "red" {
		t.Fatal("not red when there are results")
	}
}

func TestIdleQueriesStatus(t *testing.T) {
	values := make([]idleQueriesResult, 0)
	if idleQueriesStatus(values) != "green" {
		t.Fatal("not green on empty results")
	}

	values = make([]idleQueriesResult, 1)
	if idleQueriesStatus(values) != "red" {
		t.Fatal("not red when there are results")
	}
}

func TestUnusedIndexesStatus(t *testing.T) {
	values := make([]unusedIndexesResult, 0)
	if unusedIndexesStatus(values) != "green" {
		t.Fatal("not green on empty results")
	}

	values = make([]unusedIndexesResult, 1)
	if unusedIndexesStatus(values) != "red" {
		t.Fatal("not red when there are results")
	}
}

func TestBloatStatus(t *testing.T) {
	values := make([]bloatResult, 0)
	if bloatStatus(values) != "green" {
		t.Fatal("not green on empty results")
	}

	values = make([]bloatResult, 1)
	if bloatStatus(values) != "red" {
		t.Fatal("not red when there are results")
	}
}

func TestHitRateStatus(t *testing.T) {
	values := make([]hitRateResult, 0)
	if hitRateStatus(values) != "green" {
		t.Fatal("not green on empty results")
	}

	values = make([]hitRateResult, 1)
	if hitRateStatus(values) != "red" {
		t.Fatal("not red when there are results")
	}
}

func TestBlockingStatus(t *testing.T) {
	values := make([]blockingResult, 0)
	if blockingStatus(values) != "green" {
		t.Fatal("not green on empty results")
	}

	values = make([]blockingResult, 1)
	if blockingStatus(values) != "red" {
		t.Fatal("not red when there are results")
	}
}
