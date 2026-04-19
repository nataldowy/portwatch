package alert

import (
	"time"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/scanner"
)

// HistoryNotifier is a Notifier that records events to a history.Log.
type HistoryNotifier struct {
	log *history.Log
}

// NewHistoryNotifier creates a HistoryNotifier backed by the given log.
func NewHistoryNotifier(l *history.Log) *HistoryNotifier {
	return &HistoryNotifier{log: l}
}

// Notify persists the event to the history log.
func (h *HistoryNotifier) Notify(event string, p scanner.Port) error {
	return h.log.Append(history.Entry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		Port:      p.Number,
		Proto:     p.Proto,
	})
}
