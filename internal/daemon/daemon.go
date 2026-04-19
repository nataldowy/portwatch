package daemon

import (
	"context"
	"log"
	"time"

	"portwatch/internal/alert"
	"portwatch/internal/config"
	"portwatch/internal/scanner"
)

// Daemon periodically scans ports and dispatches alerts on changes.
type Daemon struct {
	cfg        *config.Config
	scanner    *scanner.Scanner
	dispatcher *alert.Dispatcher
}

// New creates a Daemon wired up with the provided config.
func New(cfg *config.Config, dispatcher *alert.Dispatcher) *Daemon {
	s := scanner.NewScanner(cfg.PortRange.Low, cfg.PortRange.High)
	return &Daemon{
		cfg:        cfg,
		scanner:    s,
		dispatcher: dispatcher,
	}
}

// Run starts the scan loop, blocking until ctx is cancelled.
func (d *Daemon) Run(ctx context.Context) error {
	log.Printf("portwatch starting: range %d-%d, interval %s",
		d.cfg.PortRange.Low, d.cfg.PortRange.High, d.cfg.Interval)

	prev, err := d.scanner.Scan()
	if err != nil {
		return err
	}
	log.Printf("initial snapshot: %d open port(s)", len(prev.Ports))

	ticker := time.NewTicker(d.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("portwatch shutting down")
			return ctx.Err()
		case <-ticker.C:
			curr, err := d.scanner.Scan()
			if err != nil {
				log.Printf("scan error: %v", err)
				continue
			}
			d.dispatcher.Dispatch(prev, curr)
			prev = curr
		}
	}
}
