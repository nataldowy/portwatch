package watch

import (
	"testing"
	"time"
)

func TestGateBlocksBelowThreshold(t *testing.T) {
	g := NewGate(3, time.Minute)
	if g.Allow("port:80") {
		t.Fatal("expected gate to block on first occurrence")
	}
	if g.Allow("port:80") {
		t.Fatal("expected gate to block on second occurrence")
	}
}

func TestGateAllowsAtThreshold(t *testing.T) {
	g := NewGate(3, time.Minute)
	g.Allow("port:80")
	g.Allow("port:80")
	if !g.Allow("port:80") {
		t.Fatal("expected gate to allow at threshold")
	}
}

func TestGateAllowsAboveThreshold(t *testing.T) {
	g := NewGate(2, time.Minute)
	g.Allow("port:443")
	g.Allow("port:443")
	if !g.Allow("port:443") {
		t.Fatal("expected gate to allow above threshold")
	}
}

func TestGateIndependentKeys(t *testing.T) {
	g := NewGate(2, time.Minute)
	g.Allow("port:80")
	g.Allow("port:80")
	// port:443 has only one occurrence — should be blocked
	if g.Allow("port:443") {
		t.Fatal("expected independent key to be blocked")
	}
}

func TestGateExpiredEventsNotCounted(t *testing.T) {
	now := time.Now()
	g := NewGate(2, 10*time.Second)
	g.now = func() time.Time { return now }

	g.Allow("port:22")
	g.Allow("port:22")

	// advance time past the window
	g.now = func() time.Time { return now.Add(15 * time.Second) }

	// previous events are now stale; this is the first fresh one
	if g.Allow("port:22") {
		t.Fatal("expected stale events to be pruned")
	}
}

func TestGateCountReflectsWindow(t *testing.T) {
	now := time.Now()
	g := NewGate(5, 30*time.Second)
	g.now = func() time.Time { return now }

	g.Allow("port:8080")
	g.Allow("port:8080")

	if got := g.Count("port:8080"); got != 2 {
		t.Fatalf("expected count 2, got %d", got)
	}
}

func TestGateResetClearsState(t *testing.T) {
	g := NewGate(2, time.Minute)
	g.Allow("port:9090")
	g.Allow("port:9090")
	g.Reset("port:9090")

	if g.Count("port:9090") != 0 {
		t.Fatal("expected count 0 after reset")
	}
	if g.Allow("port:9090") {
		t.Fatal("expected gate to block after reset")
	}
}

func TestGateDefaultsInvalidArgs(t *testing.T) {
	g := NewGate(0, 0)
	if g.threshold != 1 {
		t.Fatalf("expected threshold default 1, got %d", g.threshold)
	}
	if g.window != time.Minute {
		t.Fatalf("expected window default 1m, got %v", g.window)
	}
}
