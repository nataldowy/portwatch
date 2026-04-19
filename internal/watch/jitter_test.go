package watch

import (
	"testing"
	"time"
)

func TestJitterStaysWithinBounds(t *testing.T) {
	j := NewJitter(0.2)
	base := 100 * time.Millisecond

	for i := 0; i < 200; i++ {
		got := j.Apply(base)
		low := time.Duration(float64(base) * 0.8)
		high := time.Duration(float64(base) * 1.2)
		if got < low || got > high {
			t.Errorf("jitter out of bounds: got %v, want [%v, %v]", got, low, high)
		}
	}
}

func TestJitterDefaultsInvalidFactor(t *testing.T) {
	j := NewJitter(-1)
	if j.factor != 0.1 {
		t.Errorf("expected factor 0.1, got %v", j.factor)
	}

	j2 := NewJitter(0)
	if j2.factor != 0.1 {
		t.Errorf("expected factor 0.1 for zero, got %v", j2.factor)
	}
}

func TestJitterProducesVariance(t *testing.T) {
	j := NewJitter(0.5)
	base := 1 * time.Second
	seen := make(map[time.Duration]bool)

	for i := 0; i < 50; i++ {
		seen[j.Apply(base)] = true
	}
	if len(seen) < 5 {
		t.Errorf("expected variance in jitter output, got only %d distinct values", len(seen))
	}
}

func TestJitterZeroBaseReturnsBase(t *testing.T) {
	j := NewJitter(0.3)
	got := j.Apply(0)
	// delta will be 0, result 0 triggers fallback to base (0)
	if got != 0 {
		t.Errorf("expected 0, got %v", got)
	}
}

func TestJitterResetDoesNotPanic(t *testing.T) {
	j := NewJitter(0.1)
	j.Reset()
	got := j.Apply(50 * time.Millisecond)
	if got <= 0 {
		t.Errorf("expected positive duration after reset, got %v", got)
	}
}
