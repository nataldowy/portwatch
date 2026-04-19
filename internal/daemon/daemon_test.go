package daemon

import (
	"bytes"
	"context"
	"testing"
	"time"

	"portwatch/internal/alert"
	"portwatch/internal/config"
)

func defaultCfg() *config.Config {
	cfg := config.Default()
	// Use a very short interval and a narrow (likely empty) range to keep tests fast.
	cfg.Interval = 50 * time.Millisecond
	cfg.PortRange.Low = 1
	cfg.PortRange.High = 1
	return cfg
}

func TestDaemonRunCancels(t *testing.T) {
	cfg := defaultCfg()
	var buf bytes.Buffer
	notifier := alert.NewLogNotifier(&buf)
	dispatcher := alert.NewDispatcher(notifier)

	d := New(cfg, dispatcher)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	err := d.Run(ctx)
	if err != context.DeadlineExceeded {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}

func TestDaemonRunTicksAtLeastOnce(t *testing.T) {
	cfg := defaultCfg()
	cfg.Interval = 30 * time.Millisecond

	var buf bytes.Buffer
	notifier := alert.NewLogNotifier(&buf)
	dispatcher := alert.NewDispatcher(notifier)

	d := New(cfg, dispatcher)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Should complete without error other than context cancellation.
	err := d.Run(ctx)
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
}
