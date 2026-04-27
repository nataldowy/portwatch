package watch

import (
	"testing"
	"time"
)

func TestTallyRecordIncrementsCount(t *testing.T) {
	tl := NewTally()
	if n := tl.Record("port:80"); n != 1 {
		t.Fatalf("expected 1, got %d", n)
	}
	if n := tl.Record("port:80"); n != 2 {
		t.Fatalf("expected 2, got %d", n)
	}
}

func TestTallyGetMissingKeyReturnsNil(t *testing.T) {
	tl := NewTally()
	if e := tl.Get("missing"); e != nil {
		t.Fatalf("expected nil, got %+v", e)
	}
}

func TestTallyFirstAndLastSeen(t *testing.T) {
	tl := NewTally()
	before := time.Now()
	tl.Record("k")
	time.Sleep(2 * time.Millisecond)
	tl.Record("k")
	after := time.Now()

	e := tl.Get("k")
	if e == nil {
		t.Fatal("expected entry")
	}
	if e.FirstSeen.Before(before) || e.FirstSeen.After(after) {
		t.Errorf("FirstSeen out of range: %v", e.FirstSeen)
	}
	if !e.LastSeen.After(e.FirstSeen) {
		t.Errorf("LastSeen should be after FirstSeen")
	}
}

func TestTallySnapshotIsIndependent(t *testing.T) {
	tl := NewTally()
	tl.Record("a")
	snap := tl.Snapshot()
	tl.Record("a")
	if snap["a"].Count != 1 {
		t.Fatalf("snapshot should not reflect later mutations")
	}
}

func TestTallyResetClearsEntries(t *testing.T) {
	tl := NewTally()
	tl.Record("x")
	tl.Reset()
	if e := tl.Get("x"); e != nil {
		t.Fatalf("expected nil after reset, got %+v", e)
	}
	if n := len(tl.Snapshot()); n != 0 {
		t.Fatalf("expected empty snapshot, got %d entries", n)
	}
}

func TestTallyIndependentKeys(t *testing.T) {
	tl := NewTally()
	tl.Record("a")
	tl.Record("a")
	tl.Record("b")
	if tl.Get("a").Count != 2 {
		t.Errorf("expected count 2 for a")
	}
	if tl.Get("b").Count != 1 {
		t.Errorf("expected count 1 for b")
	}
}
