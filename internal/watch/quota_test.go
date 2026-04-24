package watch

import (
	"testing"
	"time"
)

func TestQuotaAllowsUpToMax(t *testing.T) {
	q := NewQuota(3, time.Minute)
	for i := 0; i < 3; i++ {
		if !q.Allow("k") {
			t.Fatalf("expected allow on attempt %d", i+1)
		}
	}
	if q.Allow("k") {
		t.Fatal("expected deny after max reached")
	}
}

func TestQuotaIndependentKeys(t *testing.T) {
	q := NewQuota(1, time.Minute)
	if !q.Allow("a") {
		t.Fatal("expected allow for key a")
	}
	if !q.Allow("b") {
		t.Fatal("expected allow for key b (independent)")
	}
	if q.Allow("a") {
		t.Fatal("expected deny for key a after quota exhausted")
	}
}

func TestQuotaWindowExpiry(t *testing.T) {
	q := NewQuota(1, 20*time.Millisecond)
	if !q.Allow("k") {
		t.Fatal("expected first allow")
	}
	if q.Allow("k") {
		t.Fatal("expected deny within window")
	}
	time.Sleep(30 * time.Millisecond)
	if !q.Allow("k") {
		t.Fatal("expected allow after window expiry")
	}
}

func TestQuotaRemaining(t *testing.T) {
	q := NewQuota(5, time.Minute)
	if r := q.Remaining("k"); r != 5 {
		t.Fatalf("expected 5 remaining, got %d", r)
	}
	q.Allow("k")
	q.Allow("k")
	if r := q.Remaining("k"); r != 3 {
		t.Fatalf("expected 3 remaining, got %d", r)
	}
}

func TestQuotaResetAllowsRepeat(t *testing.T) {
	q := NewQuota(1, time.Minute)
	q.Allow("k")
	if q.Allow("k") {
		t.Fatal("expected deny before reset")
	}
	q.Reset("k")
	if !q.Allow("k") {
		t.Fatal("expected allow after reset")
	}
}

func TestQuotaDefaultsInvalidParams(t *testing.T) {
	q := NewQuota(0, 0)
	if !q.Allow("k") {
		t.Fatal("expected allow with defaulted max=1")
	}
	if q.Allow("k") {
		t.Fatal("expected deny after defaulted max exhausted")
	}
}
