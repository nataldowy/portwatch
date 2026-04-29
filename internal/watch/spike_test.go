package watch

import (
	"testing"
	"time"
)

func TestSpikeNoSpikeOnFirstEvent(t *testing.T) {
	s := NewSpike(10*time.Second, 2.0)
	if s.Allow("p:80") {
		t.Fatal("expected no spike on first event")
	}
}

func TestSpikeDetectedAboveMultiplier(t *testing.T) {
	s := NewSpike(10*time.Second, 2.0)
	// Seed baseline with a small count.
	s.Allow("p:80") // baseline seeded at 1

	// Fire many events rapidly to exceed 2× baseline.
	var spiked bool
	for i := 0; i < 10; i++ {
		if s.Allow("p:80") {
			spiked = true
			break
		}
	}
	if !spiked {
		t.Fatal("expected spike to be detected after burst")
	}
}

func TestSpikeIndependentKeys(t *testing.T) {
	s := NewSpike(10*time.Second, 2.0)
	s.Allow("p:80") // seed key 1
	s.Allow("p:443") // seed key 2

	// Burst only on key 1.
	var spiked bool
	for i := 0; i < 10; i++ {
		if s.Allow("p:80") {
			spiked = true
			break
		}
	}
	if !spiked {
		t.Fatal("expected spike on p:80")
	}
	// p:443 should not be affected.
	if s.Allow("p:443") {
		t.Fatal("p:443 should not spike independently")
	}
}

func TestSpikeResetClearsState(t *testing.T) {
	s := NewSpike(10*time.Second, 2.0)
	s.Allow("p:80")
	for i := 0; i < 10; i++ {
		s.Allow("p:80")
	}
	s.Reset("p:80")
	// After reset the baseline is gone; first call seeds it again.
	if s.Allow("p:80") {
		t.Fatal("expected no spike immediately after reset")
	}
}

func TestSpikeDefaultsInvalidMultiplier(t *testing.T) {
	s := NewSpike(10*time.Second, 0)
	if s.multiplier < 1.5 {
		t.Fatalf("expected multiplier >= 1.5, got %v", s.multiplier)
	}
}

func TestSpikeDefaultsInvalidWindow(t *testing.T) {
	s := NewSpike(-1*time.Second, 2.0)
	if s.window <= 0 {
		t.Fatalf("expected positive window, got %v", s.window)
	}
}

func TestSpikeResetAllClearsAllKeys(t *testing.T) {
	s := NewSpike(10*time.Second, 2.0)
	s.Allow("p:80")
	s.Allow("p:443")
	s.ResetAll()
	if len(s.events) != 0 {
		t.Fatalf("expected empty events after ResetAll, got %d entries", len(s.events))
	}
}
