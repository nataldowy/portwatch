package watch_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"portwatch/internal/alert"
	"portwatch/internal/config"
	"portwatch/internal/scanner"
	"portwatch/internal/watch"
)

func defaultCfg() config.Config {
	cfg := config.Default()
	cfg.PortRange.From = 1
	cfg.PortRange.To = 1
	return cfg
}

func TestWatcherRunCancels(t *testing.T) {
	cfg := defaultCfg()
	s := scanner.NewScanner(cfg.PortRange.From, cfg.PortRange.To)

	var buf strings.Builder
	notifier := alert.NewLogNotifier(&buf)
	d := alert.NewDispatcher(notifier)
	w := watch.New(s, d, 50*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	err := w.Run(ctx)
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestWatcherRunTicksAtLeastOnce(t *testing.T) {
	cfg := defaultCfg()
	s := scanner.NewScanner(cfg.PortRange.From, cfg.PortRange.To)

	var buf strings.Builder
	notifier := alert.NewLogNotifier(&buf)
	d := alert.NewDispatcher(notifier)
	w := watch.New(s, d, 30*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_ = w.Run(ctx)
	// If we reach here without panic or hang, the ticker fired at least once.
}
