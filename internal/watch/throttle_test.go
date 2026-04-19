package watch

import (
	"testing"
	"time"
)

func TestThrottleAllowsFirstOccurrence(t *testing.T) {
	th := NewThrottle(1 * time.Minute)
	if !th.Allow("tcp:8080") {
		t.Fatal("expected first occurrence to be allowed")
	}
}

func TestThrottleBlocksWithinCooldown(t *testing.T) {
	th := NewThrottle(1 * time.Minute)
	th.Allow("tcp:8080")
	if th.Allow("tcp:8080") {
		t.Fatal("expected second occurrence within cooldown to be blocked")
	}
}

func TestThrottleAllowsAfterCooldown(t *testing.T) {
	th := NewThrottle(10 * time.Millisecond)
	th.Allow("tcp:9090")
	time.Sleep(20 * time.Millisecond)
	if !th.Allow("tcp:9090") {
		t.Fatal("expected occurrence after cooldown to be allowed")
	}
}

func TestThrottleIndependentKeys(t *testing.T) {
	th := NewThrottle(1 * time.Minute)
	th.Allow("tcp:8080")
	if !th.Allow("tcp:9090") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestThrottleReset(t *testing.T) {
	th := NewThrottle(1 * time.Minute)
	th.Allow("tcp:8080")
	th.Reset()
	if !th.Allow("tcp:8080") {
		t.Fatal("expected key to be allowed after reset")
	}
}
