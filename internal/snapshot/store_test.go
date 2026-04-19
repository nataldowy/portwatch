package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"portwatch/internal/scanner"
	"portwatch/internal/snapshot"
)

func makeSnapshot() scanner.Snapshot {
	return scanner.Snapshot{
		Timestamp: time.Now().Truncate(time.Second),
		Ports: []scanner.Port{
			{Port: 80, Proto: "tcp", State: "open"},
			{Port: 443, Proto: "tcp", State: "open"},
		},
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state", "snapshot.json")
	store := snapshot.NewStore(path)

	snap := makeSnapshot()
	if err := store.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	rec, ok, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !ok {
		t.Fatal("expected record to exist")
	}
	if len(rec.Snapshot.Ports) != len(snap.Ports) {
		t.Errorf("ports len: got %d, want %d", len(rec.Snapshot.Ports), len(snap.Ports))
	}
	if rec.SavedAt.IsZero() {
		t.Error("SavedAt should not be zero")
	}
}

func TestLoadMissingFile(t *testing.T) {
	store := snapshot.NewStore(filepath.Join(t.TempDir(), "no-such-file.json"))
	_, ok, err := store.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected ok=false for missing file")
	}
}

func TestSaveCreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a", "b", "c", "snap.json")
	store := snapshot.NewStore(path)
	if err := store.Save(makeSnapshot()); err != nil {
		t.Fatalf("Save with nested dirs: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}
