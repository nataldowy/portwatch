package watch

import (
	"testing"
	"time"
)

func TestCooldownAllowsFirstOccurrence(t *testing.T) {
	c := NewCooldown(5 * time.Second)
	if !c.Allow("port:8080") {
		t.Fatal("expected first occurrence to be allowed")
	}
}

func TestCooldownBlocksWithinWindow(t *testing.T) {
	c := NewCooldown(5 * time.Second)
	c.Allow("port:8080")
	if c.Allow("port:8080") {
		t.Fatal("expected second occurrence within window to be blocked")
	}
}

func TestCooldownAllowsAfterWindow(t *testing.T) {
	c := NewCooldown(10 * time.Millisecond)
	c.Allow("port:9090")
	time.Sleep(20 * time.Millisecond)
	if !c.Allow("port:9090") {
		t.Fatal("expected occurrence after window to be allowed")
	}
}

func TestCooldownIndependentKeys(t *testing.T) {
	c := NewCooldown(5 * time.Second)
	c.Allow("port:8080")
	if !c.Allow("port:9090") {
		t.Fatal("expected independent key to be allowed")
	}
}

func TestCooldownReset(t *testing.T) {
	c := NewCooldown(5 * time.Second)
	c.Allow("port:8080")
	c.Reset("port:8080")
	if !c.Allow("port:8080") {
		t.Fatal("expected allow after reset")
	}
}

func TestCooldownLen(t *testing.T) {
	c := NewCooldown(5 * time.Second)
	c.Allow("a")
	c.Allow("b")
	c.Allow("c")
	if c.Len() != 3 {
		t.Fatalf("expected 3 keys, got %d", c.Len())
	}
}
