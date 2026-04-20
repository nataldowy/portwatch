package watch

// WatchEvent names emitted by the watcher via Observer.
const (
	EventScanComplete = "scan.complete"
	EventAlertEmitted = "alert.emitted"
	EventScanError    = "scan.error"
	EventDaemonStart  = "daemon.start"
	EventDaemonStop   = "daemon.stop"
)

// ScanCompletePayload is attached to EventScanComplete events.
type ScanCompletePayload struct {
	PortsFound int
	NewPorts   int
	Closed     int
}

// AlertPayload is attached to EventAlertEmitted events.
type AlertPayload struct {
	Kind string
	Port int
	Proto string
}

// ErrorPayload is attached to EventScanError events.
type ErrorPayload struct {
	Err error
}

// ObservedWatcher wraps a watcher Observer and provides typed emit helpers
// so callers never need to remember raw event name strings.
type ObservedWatcher struct {
	obs *Observer
}

// NewObservedWatcher creates an ObservedWatcher backed by obs.
func NewObservedWatcher(obs *Observer) *ObservedWatcher {
	return &ObservedWatcher{obs: obs}
}

// ScanComplete emits a EventScanComplete event.
func (ow *ObservedWatcher) ScanComplete(p ScanCompletePayload) {
	ow.obs.Emit(EventScanComplete, p)
}

// AlertEmitted emits a EventAlertEmitted event.
func (ow *ObservedWatcher) AlertEmitted(p AlertPayload) {
	ow.obs.Emit(EventAlertEmitted, p)
}

// ScanError emits a EventScanError event.
func (ow *ObservedWatcher) ScanError(err error) {
	ow.obs.Emit(EventScanError, ErrorPayload{Err: err})
}
