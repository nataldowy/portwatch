package watch

import (
	"testing"
	"time"
)

func TestSuppressAllowsFirstOccurrence(t *testing.T) {
	s := NewSuppress(time.Second)
	if !s.Allow("port:80") {
		t.Fatal("expected first occurrence to be allowed")
	}
}

func TestSuppressBlocksWithinWindow(t *testing.T) {
	s := NewSuppress(time.Second)
	s.Allow("port:80") // prime
	if s.Allow("port:80") {
		t.Fatal("expected second occurrence within window to be suppressed")
	}
}

func TestSuppressAllowsAfterWindowExpires(t *testing.T) {
	now := time.Now()
	s := NewSuppress(50 * time.Millisecond)
	s.now = func() time.Time { return now }
	s.Allow("port:80")

	// Advance past the window.
	s.now = func() time.Time { return now.Add(100 * time.Millisecond) }
	if !s.Allow("port:80") {
		t.Fatal("expected allow after window expiry")
	}
}

func TestSuppressExtendsDeadlineOnHit(t *testing.T) {
	now := time.Now()
	s := NewSuppress(100 * time.Millisecond)
	s.now = func() time.Time { return now }
	s.Allow("port:80") // t=0, deadline=100ms

	// At t=80ms, still suppressed; deadline pushed to 180ms.
	s.now = func() time.Time { return now.Add(80 * time.Millisecond) }
	s.Allow("port:80")

	// At t=120ms, still within extended deadline.
	s.now = func() time.Time { return now.Add(120 * time.Millisecond) }
	if s.Allow("port:80") {
		t.Fatal("expected suppression to persist after deadline extension")
	}
}

func TestSuppressClearUnblocks(t *testing.T) {
	s := NewSuppress(time.Second)
	s.Allow("port:80")
	s.Clear("port:80")
	if !s.Allow("port:80") {
		t.Fatal("expected allow after explicit clear")
	}
}

func TestSuppressIndependentKeys(t *testing.T) {
	s := NewSuppress(time.Second)
	s.Allow("port:80")
	if !s.Allow("port:443") {
		t.Fatal("expected independent key to be unaffected")
	}
}

func TestSuppressResetClearsAll(t *testing.T) {
	s := NewSuppress(time.Second)
	s.Allow("port:80")
	s.Allow("port:443")
	s.Reset()
	if !s.Allow("port:80") || !s.Allow("port:443") {
		t.Fatal("expected all keys to be cleared after reset")
	}
}

func TestSuppressDefaultsInvalidWindow(t *testing.T) {
	s := NewSuppress(0)
	if s.window != 30*time.Second {
		t.Fatalf("expected default window 30s, got %v", s.window)
	}
}
