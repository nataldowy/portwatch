package report_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"portwatch/internal/history"
	"portwatch/internal/report"
)

func tempLog(t *testing.T) *history.Log {
	t.Helper()
	dir := t.TempDir()
	log, err := history.NewLog(filepath.Join(dir, "history.jsonl"))
	if err != nil {
		t.Fatalf("NewLog: %v", err)
	}
	return log
}

func TestSummariseEmpty(t *testing.T) {
	log := tempLog(t)
	s, err := report.Summarise(log)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.TotalEvents != 0 {
		t.Errorf("expected 0 events, got %d", s.TotalEvents)
	}
}

func TestSummariseCounts(t *testing.T) {
	log := tempLog(t)
	now := time.Now()
	entries := []history.Entry{
		{Timestamp: now, Kind: "new", Port: 80, Proto: "tcp"},
		{Timestamp: now, Kind: "new", Port: 443, Proto: "tcp"},
		{Timestamp: now, Kind: "closed", Port: 8080, Proto: "tcp"},
	}
	for _, e := range entries {
		if err := log.Append(e); err != nil {
			t.Fatalf("Append: %v", err)
		}
	}
	s, err := report.Summarise(log)
	if err != nil {
		t.Fatalf("Summarise: %v", err)
	}
	if s.TotalEvents != 3 {
		t.Errorf("TotalEvents: want 3, got %d", s.TotalEvents)
	}
	if s.NewPorts != 2 {
		t.Errorf("NewPorts: want 2, got %d", s.NewPorts)
	}
	if s.ClosedPorts != 1 {
		t.Errorf("ClosedPorts: want 1, got %d", s.ClosedPorts)
	}
}

func TestPrintOutput(t *testing.T) {
	s := report.Summary{TotalEvents: 5, NewPorts: 3, ClosedPorts: 2,
		FirstSeen: time.Now().Add(-time.Hour), LastSeen: time.Now()}
	var buf bytes.Buffer
	report.Print(&buf, s)
	out := buf.String()
	for _, want := range []string{"Total events", "New ports", "Closed ports", "First event", "Last event"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q", want)
		}
	}
	_ = os.Stdout
}
