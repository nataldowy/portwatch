package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempLog(t *testing.T) *Log {
	t.Helper()
	dir := t.TempDir()
	return NewLog(filepath.Join(dir, "sub", "history.jsonl"))
}

func TestAppendAndReadAll(t *testing.T) {
	l := tempLog(t)

	e1 := Entry{Timestamp: time.Now().UTC(), Event: "opened", Port: 8080, Proto: "tcp"}
	e2 := Entry{Timestamp: time.Now().UTC(), Event: "closed", Port: 22, Proto: "tcp"}

	if err := l.Append(e1); err != nil {
		t.Fatalf("append e1: %v", err)
	}
	if err := l.Append(e2); err != nil {
		t.Fatalf("append e2: %v", err)
	}

	entries, err := l.ReadAll()
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Port != 8080 || entries[0].Event != "opened" {
		t.Errorf("unexpected first entry: %+v", entries[0])
	}
	if entries[1].Port != 22 || entries[1].Event != "closed" {
		t.Errorf("unexpected second entry: %+v", entries[1])
	}
}

func TestReadAllMissingFile(t *testing.T) {
	l := NewLog(filepath.Join(t.TempDir(), "missing.jsonl"))
	entries, err := l.ReadAll()
	if err != nil {
		t.Fatalf("expected nil error for missing file, got: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil slice, got %v", entries)
	}
}

func TestAppendCreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a", "b", "c", "history.jsonl")
	l := NewLog(path)

	if err := l.Append(Entry{Timestamp: time.Now().UTC(), Event: "opened", Port: 443, Proto: "tcp"}); err != nil {
		t.Fatalf("append: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}
