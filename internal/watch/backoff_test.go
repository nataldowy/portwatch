package watch

import (
	"testing"
	"time"
)

func TestBackoffFirstFailure(t *testing.T) {
	b := NewBackoff(100*time.Millisecond, 2*time.Second)
	d := b.Record("host:22")
	if d != 100*time.Millisecond {
		t.Fatalf("expected 100ms, got %v", d)
	}
}

func TestBackoffExponentialGrowth(t *testing.T) {
	b := NewBackoff(100*time.Millisecond, 10*time.Second)
	expected := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		400 * time.Millisecond,
		800 * time.Millisecond,
	}
	for i, want := range expected {
		got := b.Record("k")
		if got != want {
			t.Fatalf("step %d: expected %v, got %v", i, want, got)
		}
	}
}

func TestBackoffCapsAtMax(t *testing.T) {
	b := NewBackoff(500*time.Millisecond, 1*time.Second)
	var last time.Duration
	for i := 0; i < 10; i++ {
		last = b.Record("x")
	}
	if last != 1*time.Second {
		t.Fatalf("expected max 1s, got %v", last)
	}
}

func TestBackoffResetClearsCount(t *testing.T) {
	b := NewBackoff(100*time.Millisecond, 5*time.Second)
	b.Record("y")
	b.Record("y")
	b.Record("y")
	b.Reset("y")
	if f := b.Failures("y"); f != 0 {
		t.Fatalf("expected 0 failures after reset, got %d", f)
	}
	d := b.Record("y")
	if d != 100*time.Millisecond {
		t.Fatalf("expected base delay after reset, got %v", d)
	}
}

func TestBackoffIndependentKeys(t *testing.T) {
	b := NewBackoff(100*time.Millisecond, 10*time.Second)
	b.Record("a")
	b.Record("a")
	d := b.Record("b")
	if d != 100*time.Millisecond {
		t.Fatalf("keys should be independent, got %v", d)
	}
}
