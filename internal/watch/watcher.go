package watch

import (
	"context"
	"time"

	"portwatch/internal/alert"
	"portwatch/internal/scanner"
)

// Watcher continuously scans ports and dispatches alerts on changes.
type Watcher struct {
	scanner    *scanner.Scanner
	dispatcher *alert.Dispatcher
	interval   time.Duration
}

// New creates a Watcher with the given scanner, dispatcher, and poll interval.
func New(s *scanner.Scanner, d *alert.Dispatcher, interval time.Duration) *Watcher {
	return &Watcher{
		scanner:    s,
		dispatcher: d,
		interval:   interval,
	}
}

// Run starts the watch loop. It blocks until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	prev, err := w.scanner.Scan()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			curr, err := w.scanner.Scan()
			if err != nil {
				return err
			}
			diff := scanner.Compare(prev, curr)
			w.dispatcher.Dispatch(diff)
			prev = curr
		}
	}
}
