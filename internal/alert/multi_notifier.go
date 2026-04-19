package alert

// MultiNotifier fans out a single event to multiple Notifier implementations.
type MultiNotifier struct {
	notifiers []Notifier
}

// NewMultiNotifier returns a MultiNotifier that dispatches to all provided notifiers.
func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{notifiers: notifiers}
}

// Notify calls Notify on every contained notifier, collecting errors.
// All notifiers are attempted even if one fails; the last non-nil error is returned.
func (m *MultiNotifier) Notify(e Event) error {
	var lastErr error
	for _, n := range m.notifiers {
		if err := n.Notify(e); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// Add appends a notifier to the fan-out list at runtime.
func (m *MultiNotifier) Add(n Notifier) {
	m.notifiers = append(m.notifiers, n)
}

// Len returns the number of registered notifiers.
func (m *MultiNotifier) Len() int {
	return len(m.notifiers)
}
