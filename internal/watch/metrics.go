package watch

import (
	"sync"
	"time"
)

// Metrics tracks runtime counters for the watcher pipeline.
type Metrics struct {
	mu           sync.RWMutex
	ScansTotal   int
	AlertsTotal  int
	ErrorsTotal  int
	LastScanAt   time.Time
	LastAlertAt  time.Time
	LastErrorAt  time.Time
	LastErrorMsg string
}

// NewMetrics returns a zeroed Metrics instance.
func NewMetrics() *Metrics {
	return &Metrics{}
}

// RecordScan increments the scan counter and updates the timestamp.
func (m *Metrics) RecordScan() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ScansTotal++
	m.LastScanAt = time.Now()
}

// RecordAlert increments the alert counter and updates the timestamp.
func (m *Metrics) RecordAlert() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.AlertsTotal++
	m.LastAlertAt = time.Now()
}

// RecordError increments the error counter and stores the message.
func (m *Metrics) RecordError(msg string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ErrorsTotal++
	m.LastErrorAt = time.Now()
	m.LastErrorMsg = msg
}

// Snapshot returns a copy of the current metrics.
func (m *Metrics) Snapshot() Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return Metrics{
		ScansTotal:   m.ScansTotal,
		AlertsTotal:  m.AlertsTotal,
		ErrorsTotal:  m.ErrorsTotal,
		LastScanAt:   m.LastScanAt,
		LastAlertAt:  m.LastAlertAt,
		LastErrorAt:  m.LastErrorAt,
		LastErrorMsg: m.LastErrorMsg,
	}
}

// Reset zeroes all counters.
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	*m = Metrics{}
}
