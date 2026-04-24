package watch

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestQuotaConcurrentAllow(t *testing.T) {
	const (
		goroutines = 50
		max        = 10
	)
	q := NewQuota(max, time.Minute)
	var allowed atomic.Int64
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			if q.Allow("shared") {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()
	if got := allowed.Load(); got != max {
		t.Fatalf("expected exactly %d allows under concurrency, got %d", max, got)
	}
}

func TestQuotaConcurrentIndependentKeys(t *testing.T) {
	const goroutines = 20
	q := NewQuota(1, time.Minute)
	var wg sync.WaitGroup
	wg.Add(goroutines)
	var denied atomic.Int64
	for i := 0; i < goroutines; i++ {
		key := itoa(i)
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			if !q.Allow(k) {
				denied.Add(1)
			}
		}(key)
		wg.Done()
	}
	wg.Wait()
	if denied.Load() != 0 {
		t.Fatalf("expected no denials for independent keys, got %d", denied.Load())
	}
}
