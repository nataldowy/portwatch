package watch

import (
	"testing"
	"time"
)

func TestCircuitAllowsWhenClosed(t *testing.T) {
	c := NewCircuit(3, time.Second)
	if !c.Allow("host:80") {
		t.Fatal("expected circuit to allow when no failures recorded")
	}
}

func TestCircuitOpensAfterThreshold(t *testing.T) {
	c := NewCircuit(3, time.Second)
	key := "host:80"
	for i := 0; i < 3; i++ {
		c.RecordFailure(key)
	}
	if c.Allow(key) {
		t.Fatal("expected circuit to be open after threshold failures")
	}
}

func TestCircuitRemainsOpenWithinWindow(t *testing.T) {
	c := NewCircuit(2, time.Minute)
	key := "host:443"
	c.RecordFailure(key)
	c.RecordFailure(key)
	if c.Allow(key) {
		t.Fatal("expected circuit to remain open within cooldown window")
	}
}

func TestCircuitHalfOpenAfterWindow(t *testing.T) {
	now := time.Now()
	c := NewCircuit(1, time.Second)
	c.now = func() time.Time { return now }
	key := "host:22"
	c.RecordFailure(key)
	// advance time past window
	c.now = func() time.Time { return now.Add(2 * time.Second) }
	if !c.Allow(key) {
		t.Fatal("expected circuit to allow probe after window expires")
	}
}

func TestCircuitSuccessResets(t *testing.T) {
	c := NewCircuit(2, time.Second)
	key := "host:8080"
	c.RecordFailure(key)
	c.RecordFailure(key)
	c.RecordSuccess(key)
	if !c.Allow(key) {
		t.Fatal("expected circuit to be closed after success")
	}
}

func TestCircuitIndependentKeys(t *testing.T) {
	c := NewCircuit(2, time.Minute)
	c.RecordFailure("a")
	c.RecordFailure("a")
	if !c.Allow("b") {
		t.Fatal("expected independent key to be unaffected")
	}
}
