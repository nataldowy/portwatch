package watch

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestSchedulerRunsTask(t *testing.T) {
	var count int32
	s := NewScheduler(10*time.Millisecond, 0, func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	s.Run(ctx)

	if atomic.LoadInt32(&count) < 2 {
		t.Fatalf("expected at least 2 task executions, got %d", count)
	}
}

func TestSchedulerRecordsLastError(t *testing.T) {
	sentinel := errors.New("scan failed")
	s := NewScheduler(10*time.Millisecond, 0, func(ctx context.Context) error {
		return sentinel
	})

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	s.Run(ctx)

	_, lastErr, _ := s.Stats()
	if !errors.Is(lastErr, sentinel) {
		t.Fatalf("expected sentinel error, got %v", lastErr)
	}
}

func TestSchedulerStatsRunCount(t *testing.T) {
	s := NewScheduler(10*time.Millisecond, 0, func(ctx context.Context) error {
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()
	s.Run(ctx)

	_, _, runs := s.Stats()
	if runs < 2 {
		t.Fatalf("expected runs >= 2, got %d", runs)
	}
}

func TestSchedulerStopsOnCancel(t *testing.T) {
	var count int32
	s := NewScheduler(10*time.Millisecond, 0, func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		s.Run(ctx)
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case <-done:
		// ok
	case <-time.After(200 * time.Millisecond):
		t.Fatal("scheduler did not stop after context cancel")
	}
}
