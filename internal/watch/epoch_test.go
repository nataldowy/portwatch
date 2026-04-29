package watch

import (
	"testing"
	"time"
)

func TestEpochStartsAtZero(t *testing.T) {
	e := NewEpoch()
	if got := e.Current(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestEpochAdvanceIncrementsGeneration(t *testing.T) {
	e := NewEpoch()
	if gen := e.Advance(); gen != 1 {
		t.Fatalf("expected 1 after first advance, got %d", gen)
	}
	if gen := e.Advance(); gen != 2 {
		t.Fatalf("expected 2 after second advance, got %d", gen)
	}
	if got := e.Current(); got != 2 {
		t.Fatalf("Current() should return 2, got %d", got)
	}
}

func TestEpochStaleDetectsOldGeneration(t *testing.T) {
	e := NewEpoch()
	old := e.Current() // 0
	e.Advance()
	if !e.Stale(old) {
		t.Fatal("expected generation 0 to be stale after advance")
	}
}

func TestEpochStaleReturnsFalseForCurrentGeneration(t *testing.T) {
	e := NewEpoch()
	e.Advance()
	cur := e.Current()
	if e.Stale(cur) {
		t.Fatal("current generation should not be stale")
	}
}

func TestEpochSinceGrowsOverTime(t *testing.T) {
	e := NewEpoch()
	time.Sleep(5 * time.Millisecond)
	if e.Since() < 5*time.Millisecond {
		t.Fatal("Since() should reflect elapsed time since creation")
	}
}

func TestEpochResetClearsGeneration(t *testing.T) {
	e := NewEpoch()
	e.Advance()
	e.Advance()
	e.Reset()
	if got := e.Current(); got != 0 {
		t.Fatalf("expected 0 after Reset, got %d", got)
	}
}

func TestEpochResetRefreshesSince(t *testing.T) {
	e := NewEpoch()
	time.Sleep(10 * time.Millisecond)
	e.Reset()
	if e.Since() >= 10*time.Millisecond {
		t.Fatal("Since() should be near-zero immediately after Reset")
	}
}
