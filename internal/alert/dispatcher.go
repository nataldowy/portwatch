package alert

import (
	"time"

	"portwatch/internal/scanner"
)

// Dispatcher converts a scanner.Diff into alert Events and forwards them
// to one or more Notifiers.
type Dispatcher struct {
	notifiers []Notifier
}

// NewDispatcher creates a Dispatcher with the supplied notifiers.
func NewDispatcher(notifiers ...Notifier) *Dispatcher {
	return &Dispatcher{notifiers: notifiers}
}

// Dispatch emits ALERT events for new ports and WARN events for closed ports.
func (d *Dispatcher) Dispatch(diff scanner.Diff) []error {
	var errs []error
	now := time.Now()

	for _, p := range diff.New {
		e := Event{
			Timestamp: now,
			Level:     LevelAlert,
			Message:   "new port detected",
			Port:      p,
		}
		for _, n := range d.notifiers {
			if err := n.Notify(e); err != nil {
				errs = append(errs, err)
			}
		}
	}

	for _, p := range diff.Closed {
		e := Event{
			Timestamp: now,
			Level:     LevelWarn,
			Message:   "port closed",
			Port:      p,
		}
		for _, n := range d.notifiers {
			if err := n.Notify(e); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errs
}
