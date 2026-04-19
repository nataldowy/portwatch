package watch

import (
	"testing"
	"time"
)

func TestRateLimitAllowsUpToMax(t *testing.T) {
	rl := NewRateLimit(time.Second, 3)
	for i := 0; i < 3; i++ {
		if !rl.Allow("k") {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
	if rl.Allow("k") {
		t.Fatal("expected block after max")
	}
}

func TestRateLimitIndependentKeys(t *testing.T) {
	rl := NewRateLimit(time.Second, 1)
	if !rl.Allow("a") {
		t.Fatal("expected allow for a")
	}
	if !rl.Allow("b") {
		t.Fatal("expected allow for b")
	}
	if rl.Allow("a") {
		t.Fatal("expected block for a")
	}
}

func TestRateLimitResetAllowsAgain(t *testing.T) {
	rl := NewRateLimit(time.Second, 1)
	rl.Allow("k")
	rl.Reset("k")
	if !rl.Allow("k") {
		t.Fatal("expected allow after reset")
	}
}

func TestRateLimitRemaining(t *testing.T) {
	rl := NewRateLimit(time.Second, 3)
	if r := rl.Remaining("k"); r != 3 {
		t.Fatalf("expected 3, got %d", r)
	}
	rl.Allow("k")
	if r := rl.Remaining("k"); r != 2 {
		t.Fatalf("expected 2, got %d", r)
	}
}

func TestRateLimitWindowExpiry(t *testing.T) {
	rl := NewRateLimit(10*time.Millisecond, 1)
	rl.Allow("k")
	if rl.Allow("k") {
		t.Fatal("expected block within window")
	}
	time.Sleep(20 * time.Millisecond)
	if !rl.Allow("k") {
		t.Fatal("expected allow after window expiry")
	}
}
