package watch

import (
	"testing"
	"time"
)

func TestThrottleGroupAllowsFirstOccurrence(t *testing.T) {
	g := NewThrottleGroup(time.Second)
	if !g.Allow("port:8080:new") {
		t.Fatal("expected first occurrence to be allowed")
	}
}

func TestThrottleGroupBlocksWithinWindow(t *testing.T) {
	g := NewThrottleGroup(time.Second)
	g.Allow("k")
	if g.Allow("k") {
		t.Fatal("expected second call within window to be blocked")
	}
}

func TestThrottleGroupAllowsAfterWindow(t *testing.T) {
	g := NewThrottleGroup(10 * time.Millisecond)
	g.Allow("k")
	time.Sleep(20 * time.Millisecond)
	if !g.Allow("k") {
		t.Fatal("expected allow after window expiry")
	}
}

func TestThrottleGroupIndependentKeys(t *testing.T) {
	g := NewThrottleGroup(time.Second)
	g.Allow("a")
	if !g.Allow("b") {
		t.Fatal("key 'b' should be independent of key 'a'")
	}
}

func TestThrottleGroupActive(t *testing.T) {
	g := NewThrottleGroup(time.Second)
	g.Allow("x")
	g.Allow("y")
	if got := g.Active(); got != 2 {
		t.Fatalf("expected 2 active, got %d", got)
	}
}

func TestThrottleGroupActiveExcludesExpired(t *testing.T) {
	g := NewThrottleGroup(10 * time.Millisecond)
	g.Allow("x")
	time.Sleep(20 * time.Millisecond)
	if got := g.Active(); got != 0 {
		t.Fatalf("expected 0 active after expiry, got %d", got)
	}
}

func TestThrottleGroupReset(t *testing.T) {
	g := NewThrottleGroup(time.Second)
	g.Allow("k")
	g.Reset()
	if !g.Allow("k") {
		t.Fatal("expected allow after reset")
	}
}

func TestThrottleGroupSetWindow(t *testing.T) {
	g := NewThrottleGroup(time.Second)
	g.Allow("k")
	g.Reset()
	g.SetWindow(10 * time.Millisecond)
	g.Allow("k")
	time.Sleep(20 * time.Millisecond)
	if !g.Allow("k") {
		t.Fatal("expected allow after updated window expiry")
	}
}

func TestThrottleGroupDefaultsInvalidWindow(t *testing.T) {
	g := NewThrottleGroup(-1)
	// Should not panic; just uses the default window.
	if g == nil {
		t.Fatal("expected non-nil ThrottleGroup")
	}
}
