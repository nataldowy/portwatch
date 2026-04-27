package watch

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestShieldConcurrentAllow(t *testing.T) {
	s := NewShield(50*time.Millisecond, 5*time.Second)
	const goroutines = 50
	var allowed int64

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			if s.Allow("shared") {
				atomic.AddInt64(&allowed, 1)
			}
		}()
	}
	wg.Wait()

	// Exactly one goroutine should have been allowed through.
	if allowed != 1 {
		t.Fatalf("expected exactly 1 allowed, got %d", allowed)
	}
}

func TestShieldConcurrentIndependentKeys(t *testing.T) {
	s := NewShield(50*time.Millisecond, 5*time.Second)
	const keys = 20
	var allowed int64

	var wg sync.WaitGroup
	wg.Add(keys)
	for i := 0; i < keys; i++ {
		key := itoa(i)
		go func(k string) {
			defer wg.Done()
			if s.Allow(k) {
				atomic.AddInt64(&allowed, 1)
			}
		}(key)
	}
	wg.Wait()

	if int(allowed) != keys {
		t.Fatalf("expected %d allowed (one per key), got %d", keys, allowed)
	}
}
