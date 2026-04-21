package watch

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSemaphoreHighConcurrencyNeverExceedsLimit(t *testing.T) {
	const limit = 5
	const workers = 50
	s := NewSemaphore(limit)
	ctx := context.Background()

	var active int64
	var violations int64
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.Acquire(ctx); err != nil {
				return
			}
			v := atomic.AddInt64(&active, 1)
			if v > int64(limit) {
				atomic.AddInt64(&violations, 1)
			}
			time.Sleep(5 * time.Millisecond)
			atomic.AddInt64(&active, -1)
			s.Release()
		}()
	}
	wg.Wait()

	if violations > 0 {
		t.Errorf("concurrency limit violated %d time(s)", violations)
	}
}

func TestSemaphoreContextCancelUnblocks(t *testing.T) {
	s := NewSemaphore(1)
	ctx := context.Background()

	// Fill the semaphore.
	if err := s.Acquire(ctx); err != nil {
		t.Fatal(err)
	}

	ctx2, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		done <- s.Acquire(ctx2)
	}()

	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err == nil {
			t.Error("expected non-nil error after context cancel")
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("Acquire did not unblock after context cancel")
	}
}
