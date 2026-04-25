package watch

import (
	"testing"
	"time"
)

func TestPressureNotHighWhenEmpty(t *testing.T) {
	p := NewPressure(3, time.Minute)
	if p.High() {
		t.Fatal("expected not high on empty pressure monitor")
	}
}

func TestPressureNotHighBelowThreshold(t *testing.T) {
	p := NewPressure(3, time.Minute)
	p.Record()
	p.Record()
	if p.High() {
		t.Fatal("expected not high below threshold")
	}
}

func TestPressureHighAtThreshold(t *testing.T) {
	p := NewPressure(3, time.Minute)
	p.Record()
	p.Record()
	p.Record()
	if !p.High() {
		t.Fatal("expected high at threshold")
	}
}

func TestPressureCountMatchesRecords(t *testing.T) {
	p := NewPressure(10, time.Minute)
	for i := 0; i < 5; i++ {
		p.Record()
	}
	if got := p.Count(); got != 5 {
		t.Fatalf("expected count 5, got %d", got)
	}
}

func TestPressureResetClearsEvents(t *testing.T) {
	p := NewPressure(3, time.Minute)
	p.Record()
	p.Record()
	p.Record()
	p.Reset()
	if p.High() {
		t.Fatal("expected not high after reset")
	}
	if got := p.Count(); got != 0 {
		t.Fatalf("expected count 0 after reset, got %d", got)
	}
}

func TestPressureWindowExpiry(t *testing.T) {
	p := NewPressure(2, 50*time.Millisecond)
	p.Record()
	p.Record()
	if !p.High() {
		t.Fatal("expected high before window expires")
	}
	time.Sleep(70 * time.Millisecond)
	if p.High() {
		t.Fatal("expected not high after window expired")
	}
	if got := p.Count(); got != 0 {
		t.Fatalf("expected count 0 after expiry, got %d", got)
	}
}

func TestPressureDefaultsInvalidThreshold(t *testing.T) {
	p := NewPressure(0, time.Minute)
	for i := 0; i < 10; i++ {
		p.Record()
	}
	if !p.High() {
		t.Fatal("expected high at default threshold of 10")
	}
}

func TestPressureDefaultsInvalidWindow(t *testing.T) {
	p := NewPressure(5, 0)
	if p.window != time.Minute {
		t.Fatalf("expected default window of 1m, got %v", p.window)
	}
}
