package watch

import (
	"testing"
	"time"
)

func TestWindowAllowsUpToMax(t *testing.T) {
	w := NewWindow(time.Minute, 3)
	for i := 0; i < 3; i++ {
		if !w.Allow("k") {
			t.Fatalf("expected allow on iteration %d", i)
		}
	}
	if w.Allow("k") {
		t.Fatal("expected deny after max")
	}
}

func TestWindowIndependentKeys(t *testing.T) {
	w := NewWindow(time.Minute, 1)
	if !w.Allow("a") {
		t.Fatal("expected allow for a")
	}
	if !w.Allow("b") {
		t.Fatal("expected allow for b")
	}
	if w.Allow("a") {
		t.Fatal("expected deny for a after max")
	}
}

func TestWindowRemaining(t *testing.T) {
	w := NewWindow(time.Minute, 5)
	if r := w.Remaining("k"); r != 5 {
		t.Fatalf("expected 5, got %d", r)
	}
	w.Allow("k")
	w.Allow("k")
	if r := w.Remaining("k"); r != 3 {
		t.Fatalf("expected 3, got %d", r)
	}
}

func TestWindowResetRestoresAllowance(t *testing.T) {
	w := NewWindow(time.Minute, 2)
	w.Allow("k")
	w.Allow("k")
	if w.Allow("k") {
		t.Fatal("expected deny before reset")
	}
	w.Reset("k")
	if !w.Allow("k") {
		t.Fatal("expected allow after reset")
	}
}

func TestWindowExpiry(t *testing.T) {
	w := NewWindow(50*time.Millisecond, 1)
	w.Allow("k")
	if w.Allow("k") {
		t.Fatal("expected deny within window")
	}
	time.Sleep(60 * time.Millisecond)
	if !w.Allow("k") {
		t.Fatal("expected allow after window expiry")
	}
}
