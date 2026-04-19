package report_test

import (
	"testing"
	"time"

	"portwatch/internal/history"
	"portwatch/internal/report"
)

var base = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func makeEntries() []history.Entry {
	return []history.Entry{
		{Timestamp: base, Kind: "new", Port: 80},
		{Timestamp: base.Add(time.Minute), Kind: "closed", Port: 8080},
		{Timestamp: base.Add(2 * time.Minute), Kind: "new", Port: 443},
	}
}

func TestFilterByKind(t *testing.T) {
	res := report.Apply(makeEntries(), report.Filter{Kind: "new"})
	if len(res) != 2 {
		t.Errorf("want 2, got %d", len(res))
	}
}

func TestFilterBySince(t *testing.T) {
	res := report.Apply(makeEntries(), report.Filter{Since: base.Add(time.Minute)})
	if len(res) != 2 {
		t.Errorf("want 2, got %d", len(res))
	}
}

func TestFilterByUntil(t *testing.T) {
	res := report.Apply(makeEntries(), report.Filter{Until: base.Add(30 * time.Second)})
	if len(res) != 1 {
		t.Errorf("want 1, got %d", len(res))
	}
}

func TestFilterNoConstraints(t *testing.T) {
	res := report.Apply(makeEntries(), report.Filter{})
	if len(res) != 3 {
		t.Errorf("want 3, got %d", len(res))
	}
}

func TestFilterCombined(t *testing.T) {
	res := report.Apply(makeEntries(), report.Filter{
		Kind:  "new",
		Since: base.Add(time.Minute),
	})
	if len(res) != 1 || res[0].Port != 443 {
		t.Errorf("unexpected result: %+v", res)
	}
}
