package watch_test

import (
	"testing"
	"time"

	"portwatch/internal/scanner"
	"portwatch/internal/watch"
)

func TestStateSetAndGet(t *testing.T) {
	var s watch.State

	snap, ts := s.Get()
	if snap != nil {
		t.Fatal("expected nil snapshot initially")
	}
	if !ts.IsZero() {
		t.Fatal("expected zero time initially")
	}

	sc := scanner.NewScanner(1, 1)
	got, err := sc.Scan()
	if err != nil {
		t.Fatalf("scan: %v", err)
	}

	before := time.Now()
	s.Set(got)
	after := time.Now()

	snap2, ts2 := s.Get()
	if snap2 == nil {
		t.Fatal("expected non-nil snapshot after Set")
	}
	if ts2.Before(before) || ts2.After(after) {
		t.Errorf("timestamp %v out of expected range [%v, %v]", ts2, before, after)
	}
}

func TestStateAge(t *testing.T) {
	var s watch.State

	if s.Age() != 0 {
		t.Fatal("expected zero age before first Set")
	}

	sc := scanner.NewScanner(1, 1)
	snap, _ := sc.Scan()
	s.Set(snap)

	time.Sleep(10 * time.Millisecond)
	if s.Age() < 10*time.Millisecond {
		t.Error("expected age to be at least 10ms")
	}
}
