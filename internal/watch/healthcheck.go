package watch

import (
	"sync"
	"time"
)

// HealthStatus represents the current health of a monitored component.
type HealthStatus struct {
	Healthy   bool
	LastCheck time.Time
	Failures  int
	Message   string
}

// HealthCheck tracks component health based on reported successes and failures.
type HealthCheck struct {
	mu       sync.RWMutex
	states   map[string]*HealthStatus
	threshold int
}

// NewHealthCheck creates a HealthCheck where a component is considered
// unhealthy after threshold consecutive failures.
func NewHealthCheck(threshold int) *HealthCheck {
	if threshold <= 0 {
		threshold = 3
	}
	return &HealthCheck{
		states:    make(map[string]*HealthStatus),
		threshold: threshold,
	}
}

func (h *HealthCheck) getOrCreate(key string) *HealthStatus {
	if s, ok := h.states[key]; ok {
		return s
	}
	s := &HealthStatus{Healthy: true}
	h.states[key] = s
	return s
}

// RecordSuccess marks a successful check for the given key.
func (h *HealthCheck) RecordSuccess(key string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	s := h.getOrCreate(key)
	s.Failures = 0
	s.Healthy = true
	s.LastCheck = time.Now()
	s.Message = ""
}

// RecordFailure records a failure for the given key with an optional message.
func (h *HealthCheck) RecordFailure(key, message string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	s := h.getOrCreate(key)
	s.Failures++
	s.LastCheck = time.Now()
	s.Message = message
	if s.Failures >= h.threshold {
		s.Healthy = false
	}
}

// Status returns the current HealthStatus for the given key.
func (h *HealthCheck) Status(key string) HealthStatus {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if s, ok := h.states[key]; ok {
		return *s
	}
	return HealthStatus{Healthy: true}
}

// Reset clears the health state for the given key.
func (h *HealthCheck) Reset(key string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.states, key)
}
