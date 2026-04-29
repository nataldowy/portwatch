package watch

import (
	"testing"
	"time"
)

func TestStrideNotHighBelowThreshold(t *testing.T) {
	s := NewStride(3, time.Second)
	if s.Record("p:80") {
		t.Fatal("expected false with only 1 active key, threshold=3")
	}
	if s.Record("p:443") {
		t.Fatal("expected false with only 2 active keys, threshold=3")
	}
}

func TestStrideFiresAtThreshold(t *testing.T) {
	s := NewStride(3, time.Second)
	s.Record("p:80")
	s.Record("p:443")
	if !s.Record("p:8080") {
		t.Fatal("expected true when 3 unique keys active, threshold=3")
	}
}

func TestStrideFiresAboveThreshold(t *testing.T) {
	s := NewStride(2, time.Second)
	s.Record("p:80")
	if !s.Record("p:443") {
		t.Fatal("expected true at threshold=2 with 2 keys")
	}
	if !s.Record("p:8080") {
		t.Fatal("expected true above threshold")
	}
}

func TestStrideExpiredEventsNotCounted(t *testing.T) {
	now := time.Now()
	s := NewStride(2, 100*time.Millisecond)
	s.now = func() time.Time { return now }

	s.Record("p:80")
	s.Record("p:443")

	// Advance time past the window.
	s.now = func() time.Time { return now.Add(200 * time.Millisecond) }

	if s.Record("p:9090") {
		t.Fatal("old events should have expired; only 1 active key expected")
	}
}

func TestStrideActiveKeys(t *testing.T) {
	s := NewStride(5, time.Second)
	s.Record("a")
	s.Record("b")
	s.Record("a") // duplicate key, still counts as 1
	if got := s.ActiveKeys(); got != 2 {
		t.Fatalf("expected 2 active keys, got %d", got)
	}
}

func TestStrideResetClearsState(t *testing.T) {
	s := NewStride(2, time.Second)
	s.Record("p:80")
	s.Record("p:443")
	s.Reset()
	if s.ActiveKeys() != 0 {
		t.Fatal("expected 0 active keys after Reset")
	}
	if s.Record("p:80") {
		t.Fatal("expected false after reset with threshold=2")
	}
}

func TestStrideDefaultsInvalidParams(t *testing.T) {
	s := NewStride(0, -1)
	if s.threshold != 1 {
		t.Fatalf("threshold should be clamped to 1, got %d", s.threshold)
	}
	if s.window <= 0 {
		t.Fatalf("window should default to positive value, got %v", s.window)
	}
}
