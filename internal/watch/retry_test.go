package watch

import (
	"testing"
	"time"
)

func TestRetryAllowsUpToMax(t *testing.T) {
	r := NewRetry(3, 0)
	for i := 0; i < 3; i++ {
		if !r.Allow("k") {
			t.Fatalf("expected allow on attempt %d", i+1)
		}
	}
	if r.Allow("k") {
		t.Fatal("expected deny after max tries")
	}
}

func TestRetryEnforcesDelay(t *testing.T) {
	now := time.Now()
	r := NewRetry(5, 10*time.Minute)
	r.now = func() time.Time { return now }

	if !r.Allow("k") {
		t.Fatal("expected first allow")
	}
	if r.Allow("k") {
		t.Fatal("expected deny within delay")
	}
	r.now = func() time.Time { return now.Add(11 * time.Minute) }
	if !r.Allow("k") {
		t.Fatal("expected allow after delay")
	}
}

func TestRetryRemaining(t *testing.T) {
	r := NewRetry(3, 0)
	if r.Remaining("k") != 3 {
		t.Fatalf("expected 3 remaining, got %d", r.Remaining("k"))
	}
	r.Allow("k")
	if r.Remaining("k") != 2 {
		t.Fatalf("expected 2 remaining, got %d", r.Remaining("k"))
	}
}

func TestRetryResetRestoresAttempts(t *testing.T) {
	r := NewRetry(1, 0)
	r.Allow("k")
	if r.Allow("k") {
		t.Fatal("expected deny after max")
	}
	r.Reset("k")
	if !r.Allow("k") {
		t.Fatal("expected allow after reset")
	}
}

func TestRetryIndependentKeys(t *testing.T) {
	r := NewRetry(1, 0)
	r.Allow("a")
	if !r.Allow("b") {
		t.Fatal("key b should be independent of key a")
	}
}
