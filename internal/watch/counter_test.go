package watch

import (
	"sync"
	"testing"
)

func TestCounterIncAndGet(t *testing.T) {
	c := NewCounter()
	c.Inc("scans")
	c.Inc("scans")
	if got := c.Get("scans"); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestCounterAddPositive(t *testing.T) {
	c := NewCounter()
	c.Add("alerts", 5)
	if got := c.Get("alerts"); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestCounterAddIgnoresNonPositive(t *testing.T) {
	c := NewCounter()
	c.Add("errors", -3)
	c.Add("errors", 0)
	if got := c.Get("errors"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestCounterGetMissingKeyReturnsZero(t *testing.T) {
	c := NewCounter()
	if got := c.Get("nonexistent"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestCounterSnapshot(t *testing.T) {
	c := NewCounter()
	c.Inc("a")
	c.Add("b", 3)
	snap := c.Snapshot()
	if snap["a"] != 1 || snap["b"] != 3 {
		t.Fatalf("unexpected snapshot: %v", snap)
	}
	// Mutating the snapshot must not affect the counter.
	snap["a"] = 99
	if c.Get("a") != 1 {
		t.Fatal("snapshot mutation affected counter")
	}
}

func TestCounterReset(t *testing.T) {
	c := NewCounter()
	c.Inc("x")
	c.Reset()
	if got := c.Get("x"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestCounterConcurrentInc(t *testing.T) {
	c := NewCounter()
	var wg sync.WaitGroup
	const goroutines = 50
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Inc("concurrent")
		}()
	}
	wg.Wait()
	if got := c.Get("concurrent"); got != goroutines {
		t.Fatalf("expected %d, got %d", goroutines, got)
	}
}
