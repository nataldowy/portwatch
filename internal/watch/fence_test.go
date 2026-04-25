package watch

import (
	"testing"
	"time"
)

func TestFenceAllowsFirstOccurrence(t *testing.T) {
	f := NewFence(100 * time.Millisecond)
	if !f.Allow("p80") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestFenceBlocksWithinGap(t *testing.T) {
	now := time.Now()
	f := NewFence(500 * time.Millisecond)
	f.now = func() time.Time { return now }

	f.Allow("p80")
	f.now = func() time.Time { return now.Add(100 * time.Millisecond) }

	if f.Allow("p80") {
		t.Fatal("expected second call within gap to be blocked")
	}
	if got := f.Blocked("p80"); got != 1 {
		t.Fatalf("expected blocked=1, got %d", got)
	}
}

func TestFenceAllowsAfterGap(t *testing.T) {
	now := time.Now()
	f := NewFence(200 * time.Millisecond)
	f.now = func() time.Time { return now }

	f.Allow("p443")
	f.now = func() time.Time { return now.Add(300 * time.Millisecond) }

	if !f.Allow("p443") {
		t.Fatal("expected call after gap to be allowed")
	}
}

func TestFenceIndependentKeys(t *testing.T) {
	f := NewFence(1 * time.Second)
	if !f.Allow("p80") {
		t.Fatal("p80 first call should be allowed")
	}
	if !f.Allow("p443") {
		t.Fatal("p443 first call should be allowed independently")
	}
}

func TestFenceResetClearsState(t *testing.T) {
	now := time.Now()
	f := NewFence(1 * time.Second)
	f.now = func() time.Time { return now }

	f.Allow("p22")
	f.Allow("p22") // blocked

	f.Reset()

	if f.Blocked("p22") != 0 {
		t.Fatal("expected blocked counter to be cleared after Reset")
	}
	if !f.Allow("p22") {
		t.Fatal("expected allow after Reset")
	}
}

func TestFenceDefaultsInvalidGap(t *testing.T) {
	f := NewFence(0)
	if f.gap != time.Second {
		t.Fatalf("expected default gap=1s, got %v", f.gap)
	}
}
