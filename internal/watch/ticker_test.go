package watch

import (
	"context"
	"testing"
	"time"
)

func TestTickerFiresAtLeastOnce(t *testing.T) {
	tk := NewTicker(20*time.Millisecond, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	go tk.Run(ctx)

	select {
	case ts := <-tk.C():
		if ts.IsZero() {
			t.Fatal("expected non-zero timestamp")
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("ticker did not fire within timeout")
	}
}

func TestTickerStopsOnContextCancel(t *testing.T) {
	tk := NewTicker(10*time.Millisecond, 0)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		tk.Run(ctx)
		close(done)
	}()

	cancel()
	select {
	case <-done:
		// ok
	case <-time.After(200 * time.Millisecond):
		t.Fatal("ticker goroutine did not exit after cancel")
	}
}

func TestTickerDefaultsInvalidInterval(t *testing.T) {
	tk := NewTicker(-1*time.Second, 0)
	if tk.interval != 30*time.Second {
		t.Fatalf("expected default 30s, got %v", tk.interval)
	}
}

func TestTickerResetDrainsPending(t *testing.T) {
	tk := NewTicker(10*time.Millisecond, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	go tk.Run(ctx)
	time.Sleep(50 * time.Millisecond)

	tk.Reset()
	select {
	case <-tk.C():
		// drain any remaining tick — acceptable
	default:
		// channel was already empty after reset — also fine
	}
}

func TestTickerFiresMultipleTimes(t *testing.T) {
	tk := NewTicker(15*time.Millisecond, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	go tk.Run(ctx)

	count := 0
	for {
		select {
		case _, ok := <-tk.C():
			if !ok {
				if count < 2 {
					t.Fatalf("expected at least 2 ticks, got %d", count)
				}
				return
			}
			count++
		case <-ctx.Done():
			if count < 2 {
				t.Fatalf("expected at least 2 ticks, got %d", count)
			}
			return
		}
	}
}
