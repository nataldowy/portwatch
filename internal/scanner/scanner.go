package scanner

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state of a single port.
type PortState struct {
	Port     int
	Protocol string
	Open     bool
}

// Snapshot holds all open ports at a point in time.
type Snapshot struct {
	Ports     []PortState
	Timestamp time.Time
}

// Scanner scans a range of ports on localhost.
type Scanner struct {
	StartPort int
	EndPort   int
	Timeout   time.Duration
}

// NewScanner creates a Scanner with sensible defaults.
func NewScanner(start, end int, timeout time.Duration) *Scanner {
	return &Scanner{StartPort: start, EndPort: end, Timeout: timeout}
}

// Scan probes each port and returns a Snapshot.
func (s *Scanner) Scan() Snapshot {
	var ports []PortState
	for port := s.StartPort; port <= s.EndPort; port++ {
		address := fmt.Sprintf("127.0.0.1:%d", port)
		conn, err := net.DialTimeout("tcp", address, s.Timeout)
		if err == nil {
			conn.Close()
			ports = append(ports, PortState{Port: port, Protocol: "tcp", Open: true})
		}
	}
	return Snapshot{Ports: ports, Timestamp: time.Now()}
}
