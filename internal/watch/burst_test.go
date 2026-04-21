package watch

import (
	"testing"
	"time"
)

func TestBurstAllowsUpToCap(t *testing.T) {
	b := NewBurst(3, time.Second)
	for i := 0; i < 3; i++ {
		if !b.Allow("k") {
			t.Fatalf("expected Allow on iteration %d", i)
		}
	}
	if b.Allow("k") {
		t.Fatal("expected suppression after cap exceeded")
	}
}

func TestBurstIndependentKeys(t *testing.T) {
	b := NewBurst(1, time.Second)
	if !b.Allow("a") {
		t.Fatal("expected first allow for key a")
	}
	if !b.Allow("b") {
		t.Fatal("expected first allow for key b")
	}
	if b.Allow("a") {
		t.Fatal("expected suppression for key a")
	}
}

func TestBurstAllowsAfterWindowExpiry(t *testing.T) {
	b := NewBurst(1, 20*time.Millisecond)
	if !b.Allow("k") {
		t.Fatal("expected first allow")
	}
	if b.Allow("k") {
		t.Fatal("expected suppression within window")
	}
	time.Sleep(30 * time.Millisecond)
	if !b.Allow("k") {
		t.Fatal("expected allow after window expiry")
	}
}

func TestBurstRemaining(t *testing.T) {
	b := NewBurst(3, time.Second)
	if got := b.Remaining("k"); got != 3 {
		t.Fatalf("expected 3 remaining, got %d", got)
	}
	b.Allow("k")
	b.Allow("k")
	if got := b.Remaining("k"); got != 1 {
		t.Fatalf("expected 1 remaining, got %d", got)
	}
}

func TestBurstResetAllowsRepeat(t *testing.T) {
	b := NewBurst(1, time.Second)
	b.Allow("k")
	if b.Allow("k") {
		t.Fatal("expected suppression before reset")
	}
	b.Reset("k")
	if !b.Allow("k") {
		t.Fatal("expected allow after reset")
	}
}

func TestBurstDefaultsInvalidParams(t *testing.T) {
	b := NewBurst(0, 0)
	if !b.Allow("k") {
		t.Fatal("expected first allow with defaulted params")
	}
	if b.Allow("k") {
		t.Fatal("expected suppression at cap=1")
	}
}
