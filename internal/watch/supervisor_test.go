package watch

import (
	"context"
	"errors"
	"testing"
	"time"
)

var errBoom = errors.New("boom")

func TestSupervisorSuccessNoRetry(t *testing.T) {
	s := NewSupervisor(SupervisorConfig{Name: "t", MaxRetries: 3, Delay: 0})
	calls := 0
	err := s.Run(context.Background(), func(_ context.Context) error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestSupervisorRestartsOnError(t *testing.T) {
	s := NewSupervisor(SupervisorConfig{Name: "t", MaxRetries: 3, Delay: 0})
	calls := 0
	err := s.Run(context.Background(), func(_ context.Context) error {
		calls++
		if calls < 3 {
			return errBoom
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestSupervisorGivesUpAfterMaxRetries(t *testing.T) {
	s := NewSupervisor(SupervisorConfig{Name: "t", MaxRetries: 2, Delay: 0})
	calls := 0
	err := s.Run(context.Background(), func(_ context.Context) error {
		calls++
		return errBoom
	})
	if !errors.Is(err, errBoom) {
		t.Fatalf("expected errBoom, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestSupervisorRespectsContextCancel(t *testing.T) {
	s := NewSupervisor(SupervisorConfig{Name: "t", MaxRetries: 10, Delay: 50 * time.Millisecond})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
	err := s.Run(ctx, func(_ context.Context) error {
		return errBoom
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
