package scanner

import (
	"net"
	"testing"
	"time"
)

func TestScanDetectsOpenPort(t *testing.T) {
	// Start a temporary listener on a random port.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port
	s := NewScanner(port, port, 500*time.Millisecond)
	snap := s.Scan()

	if len(snap.Ports) != 1 {
		t.Fatalf("expected 1 open port, got %d", len(snap.Ports))
	}
	if snap.Ports[0].Port != port {
		t.Errorf("expected port %d, got %d", port, snap.Ports[0].Port)
	}
	if !snap.Ports[0].Open {
		t.Error("expected port to be open")
	}
}

func TestScanReturnsEmptyForClosedRange(t *testing.T) {
	// Port 1 is almost certainly closed in test environments.
	s := NewScanner(1, 1, 100*time.Millisecond)
	snap := s.Scan()
	if len(snap.Ports) != 0 {
		t.Errorf("expected no open ports, got %d", len(snap.Ports))
	}
}

func TestSnapshotTimestamp(t *testing.T) {
	s := NewScanner(1, 1, 50*time.Millisecond)
	before := time.Now()
	snap := s.Scan()
	after := time.Now()
	if snap.Timestamp.Before(before) || snap.Timestamp.After(after) {
		t.Error("snapshot timestamp is outside expected range")
	}
}
