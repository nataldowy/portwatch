package watch

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGateConcurrentAllow(t *testing.T) {
	const goroutines = 20
	const threshold = 5

	g := NewGate(threshold, time.Minute)
	var allowed atomic.Int64
	var wg sync.WaitGroup

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			if g.Allow("port:80") {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()

	// At least goroutines-threshold calls should have been allowed
	// (all calls at or above threshold count)
	if allowed.Load() == 0 {
		t.Fatal("expected at least one allowed event under concurrency")
	}
	if got := g.Count("port:80"); got != goroutines {
		t.Fatalf("expected count %d, got %d", goroutines, got)
	}
}

func TestGateConcurrentIndependentKeys(t *testing.T) {
	const goroutines = 10
	const threshold = 3

	g := NewGate(threshold, time.Minute)
	var wg sync.WaitGroup

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		key := fmt.Sprintf("port:%d", 9000+i)
		go func(k string) {
			defer wg.Done()
			// each key gets exactly one call — should never be allowed
			if g.Allow(k) {
				t.Errorf("key %s should not pass gate with one occurrence", k)
			}
		}(key)
	}
	wg.Wait()
}
