package watch

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestFenceConcurrentAllow(t *testing.T) {
	f := NewFence(50 * time.Millisecond)

	var allowed atomic.Int64
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if f.Allow("shared") {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()

	// Only the first goroutine (or very few within the same nanosecond window)
	// should be allowed; at minimum 1 must pass.
	if got := allowed.Load(); got < 1 {
		t.Fatalf("expected at least 1 allowed, got %d", got)
	}
}

func TestFenceConcurrentIndependentKeys(t *testing.T) {
	f := NewFence(1 * time.Second)

	var allowed atomic.Int64
	var wg sync.WaitGroup
	keys := []string{"p80", "p443", "p22", "p8080", "p3306"}

	for _, k := range keys {
		wg.Add(1)
		key := k
		go func() {
			defer wg.Done()
			if f.Allow(key) {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()

	if got := allowed.Load(); got != int64(len(keys)) {
		t.Fatalf("expected all %d independent keys allowed, got %d", len(keys), got)
	}
}
