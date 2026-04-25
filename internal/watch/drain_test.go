package watch

import (
	"testing"
	"time"
)

func TestDrainAddAndFlush(t *testing.T) {
	d := NewDrain(time.Minute, 10)
	now := time.Now()

	d.Add("tcp:80", now)
	d.Add("tcp:80", now.Add(time.Second))

	events := d.Flush("tcp:80", now.Add(2*time.Second))
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}

func TestDrainFlushClearsBuffer(t *testing.T) {
	d := NewDrain(time.Minute, 10)
	now := time.Now()
	d.Add("tcp:80", now)
	d.Flush("tcp:80", now)

	if d.Count("tcp:80") != 0 {
		t.Fatal("expected buffer to be empty after flush")
	}
}

func TestDrainDiscardsExpiredEvents(t *testing.T) {
	d := NewDrain(5*time.Second, 10)
	now := time.Now()

	old := now.Add(-10 * time.Second)
	d.Add("tcp:443", old)
	d.Add("tcp:443", now)

	events := d.Flush("tcp:443", now)
	if len(events) != 1 {
		t.Fatalf("expected 1 event after expiry filter, got %d", len(events))
	}
}

func TestDrainAddReturnsTrueAtMax(t *testing.T) {
	d := NewDrain(time.Minute, 3)
	now := time.Now()

	if d.Add("k", now) {
		t.Fatal("should not be full after first add")
	}
	d.Add("k", now)
	if !d.Add("k", now) {
		t.Fatal("expected true when max reached")
	}
}

func TestDrainIndependentKeys(t *testing.T) {
	d := NewDrain(time.Minute, 10)
	now := time.Now()
	d.Add("a", now)
	d.Add("b", now)
	d.Add("b", now)

	if d.Count("a") != 1 {
		t.Fatalf("expected 1 for key a, got %d", d.Count("a"))
	}
	if d.Count("b") != 2 {
		t.Fatalf("expected 2 for key b, got %d", d.Count("b"))
	}
}

func TestDrainResetClearsAll(t *testing.T) {
	d := NewDrain(time.Minute, 10)
	now := time.Now()
	d.Add("a", now)
	d.Add("b", now)
	d.Reset()

	if d.Count("a") != 0 || d.Count("b") != 0 {
		t.Fatal("expected all keys cleared after reset")
	}
}

func TestDrainDefaultsInvalidMax(t *testing.T) {
	d := NewDrain(time.Minute, 0)
	now := time.Now()
	for i := 0; i < 64; i++ {
		d.Add("x", now)
	}
	if d.Count("x") != 64 {
		t.Fatalf("expected 64 buffered, got %d", d.Count("x"))
	}
}
