package watch

import (
	"testing"
	"time"
)

func TestShieldAllowsFirstOccurrence(t *testing.T) {
	s := NewShield(100*time.Millisecond, 0)
	if !s.Allow("p:80") {
		t.Fatal("expected first Allow to return true")
	}
}

func TestShieldBlocksWithinWindow(t *testing.T) {
	s := NewShield(200*time.Millisecond, 0)
	s.Allow("p:80")
	if s.Allow("p:80") {
		t.Fatal("expected second Allow within window to return false")
	}
}

func TestShieldAllowsAfterWindow(t *testing.T) {
	s := NewShield(30*time.Millisecond, 0)
	s.Allow("p:80")
	time.Sleep(40 * time.Millisecond)
	if !s.Allow("p:80") {
		t.Fatal("expected Allow after window expiry to return true")
	}
}

func TestShieldExponentialGrowth(t *testing.T) {
	base := 20 * time.Millisecond
	s := NewShield(base, 10*time.Second)

	// First allow — window = base (20 ms)
	if !s.Allow("k") {
		t.Fatal("first allow should succeed")
	}
	// Wait past base window so second allow fires; window becomes 2×base = 40 ms
	time.Sleep(base + 5*time.Millisecond)
	if !s.Allow("k") {
		t.Fatal("second allow should succeed after first window")
	}
	// Immediately after second allow the window is 40 ms, so this should block
	if s.Allow("k") {
		t.Fatal("third allow should be blocked within doubled window")
	}
}

func TestShieldCapsAtMaxWindow(t *testing.T) {
	s := NewShield(10*time.Millisecond, 15*time.Millisecond)
	for i := 0; i < 6; i++ {
		time.Sleep(20 * time.Millisecond)
		s.Allow("k")
	}
	// After many doublings the window should be capped; sleep 20 ms and allow.
	time.Sleep(20 * time.Millisecond)
	if !s.Allow("k") {
		t.Fatal("allow should succeed after capped window expires")
	}
}

func TestShieldIndependentKeys(t *testing.T) {
	s := NewShield(200*time.Millisecond, 0)
	s.Allow("a")
	if !s.Allow("b") {
		t.Fatal("key b should be independent of key a")
	}
}

func TestShieldHitsTracked(t *testing.T) {
	s := NewShield(20*time.Millisecond, 0)
	s.Allow("k")
	time.Sleep(25 * time.Millisecond)
	s.Allow("k")
	if s.Hits("k") != 2 {
		t.Fatalf("expected 2 hits, got %d", s.Hits("k"))
	}
}

func TestShieldResetAllowsRepeat(t *testing.T) {
	s := NewShield(500*time.Millisecond, 0)
	s.Allow("k")
	s.Reset("k")
	if !s.Allow("k") {
		t.Fatal("expected Allow to succeed after Reset")
	}
	if s.Hits("k") != 1 {
		t.Fatalf("expected hits to restart at 1 after Reset, got %d", s.Hits("k"))
	}
}
