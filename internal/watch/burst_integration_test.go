package watch

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBurstConcurrentAllows(t *testing.T) {
	const cap = 10
	const goroutines = 50
	b := NewBurst(cap, time.Second)

	var allowed atomic.Int64
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			if b.Allow("shared") {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()

	if got := allowed.Load(); got > cap {
		t.Fatalf("concurrent allows exceeded cap: got %d, cap %d", got, cap)
	}
}

func TestBurstConcurrentIndependentKeys(t *testing.T) {
	const cap = 1
	const keys = 20
	b := NewBurst(cap, time.Second)

	var wg sync.WaitGroup
	wg.Add(keys)
	allowed := make([]atomic.Int64, keys)
	for i := 0; i < keys; i++ {
		go func(idx int) {
			defer wg.Done()
			key := string(rune('a' + idx))
			for j := 0; j < 5; j++ {
				if b.Allow(key) {
					allowed[idx].Add(1)
				}
			}
		}(i)
	}
	wg.Wait()

	for i := 0; i < keys; i++ {
		if got := allowed[i].Load(); got > int64(cap) {
			t.Fatalf("key %d exceeded cap: got %d", i, got)
		}
	}
}
