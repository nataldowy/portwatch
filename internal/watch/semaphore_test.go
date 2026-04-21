package watch

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestSemaphoreAcquireAndRelease(t *testing.T) {
	s := NewSemaphore(2)
	ctx := context.Background()

	if err := s.Acquire(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 1 {
		t.Errorf("expected Len 1, got %d", s.Len())
	}
	s.Release()
	if s.Len() != 0 {
		t.Errorf("expected Len 0 after release, got %d", s.Len())
	}
}

func TestSemaphoreDefaultsToOne(t *testing.T) {
	s := NewSemaphore(0)
	if s.Cap() != 1 {
		t.Errorf("expected cap 1 for invalid n, got %d", s.Cap())
	}
}

func TestSemaphoreBlocksAtCapacity(t *testing.T) {
	s := NewSemaphore(1)
	ctx := context.Background()

	if err := s.Acquire(ctx); err != nil {
		t.Fatal(err)
	}

	ctx2, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := s.Acquire(ctx2)
	if err == nil {
		t.Error("expected error when semaphore is full, got nil")
	}
}

func TestSemaphoreConcurrentAcquire(t *testing.T) {
	const limit = 3
	const goroutines = 9
	s := NewSemaphore(limit)
	ctx := context.Background()

	var mu sync.Mutex
	peak := 0
	current := 0
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.Acquire(ctx)
			mu.Lock()
			current++
			if current > peak {
				peak = current
			}
			mu.Unlock()
			time.Sleep(10 * time.Millisecond)
			mu.Lock()
			current--
			mu.Unlock()
			s.Release()
		}()
	}
	wg.Wait()

	if peak > limit {
		t.Errorf("peak concurrency %d exceeded limit %d", peak, limit)
	}
}

func TestSemaphoreResetDrainsSlots(t *testing.T) {
	s := NewSemaphore(3)
	ctx := context.Background()
	_ = s.Acquire(ctx)
	_ = s.Acquire(ctx)
	s.Reset()
	if s.Len() != 0 {
		t.Errorf("expected Len 0 after Reset, got %d", s.Len())
	}
}

func TestSemaphoreReleasePanicsWithoutAcquire(t *testing.T) {
	s := NewSemaphore(2)
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on Release without Acquire")
		}
	}()
	s.Release()
}
