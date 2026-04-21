package watch

import (
	"net"
	"sync"
	"time"
)

// ProbeResult holds the outcome of a single port probe.
type ProbeResult struct {
	Host    string
	Port    int
	Open    bool
	Latency time.Duration
	Err     error
}

// Probe performs active TCP connectivity checks against a list of ports
// and returns results. It is safe for concurrent use.
type Probe struct {
	mu      sync.Mutex
	timeout time.Duration
}

// NewProbe creates a Probe with the given dial timeout.
// If timeout is zero or negative the default of 2 seconds is used.
func NewProbe(timeout time.Duration) *Probe {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	return &Probe{timeout: timeout}
}

// Check dials host:port and returns a ProbeResult.
func (p *Probe) Check(host string, port int) ProbeResult {
	p.mu.Lock()
	timeout := p.timeout
	p.mu.Unlock()

	addr := net.JoinHostPort(host, itoa(port))
	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, timeout)
	latency := time.Since(start)

	if err != nil {
		return ProbeResult{Host: host, Port: port, Open: false, Latency: latency, Err: err}
	}
	_ = conn.Close()
	return ProbeResult{Host: host, Port: port, Open: true, Latency: latency}
}

// CheckAll probes each port in the list concurrently and returns all results.
func (p *Probe) CheckAll(host string, ports []int) []ProbeResult {
	results := make([]ProbeResult, len(ports))
	var wg sync.WaitGroup
	for i, port := range ports {
		wg.Add(1)
		go func(idx, pt int) {
			defer wg.Done()
			results[idx] = p.Check(host, pt)
		}(i, port)
	}
	wg.Wait()
	return results
}

// SetTimeout updates the dial timeout used for future checks.
func (p *Probe) SetTimeout(d time.Duration) {
	if d <= 0 {
		return
	}
	p.mu.Lock()
	p.timeout = d
	p.mu.Unlock()
}
