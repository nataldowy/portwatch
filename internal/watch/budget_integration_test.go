package watch

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBudgetConcurrentAllow(t *testing.T) {
	const goroutines = 50
	const allowance = 20
	b := NewBudget(allowance, time.Minute)
	var allowed atomic.Int64
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			if b.Allow("shared", 1) {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()
	if got := allowed.Load(); got != allowance {
		t.Fatalf("expected exactly %d allowed, got %d", allowance, got)
	}
}

func TestBudgetConcurrentIndependentKeys(t *testing.T) {
	const keys = 20
	b := NewBudget(1, time.Minute)
	var wg sync.WaitGroup
	wg.Add(keys)
	results := make([]bool, keys)
	for i := 0; i < keys; i++ {
		i := i
		go func() {
			defer wg.Done()
			results[i] = b.Allow(itoa(i), 1)
		}()
	}
	wg.Wait()
	for i, ok := range results {
		if !ok {
			t.Errorf("key %d should have been allowed", i)
		}
	}
}
