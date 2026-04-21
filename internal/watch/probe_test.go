package watch

import (
	"net"
	"strconv"
	"testing"
	"time"
)

// startTCPListener opens a local TCP listener on a random port and returns
// the listener and the chosen port number.
func startTCPListener(t *testing.T) (*net.TCPListener, int) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port, _ := strconv.Atoi(ln.Addr().(*net.TCPAddr).Port.Error())
	// Use TCPAddr directly
	port = ln.Addr().(*net.TCPAddr).Port
	return ln.(*net.TCPListener), port
}

func TestProbeOpenPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port

	p := NewProbe(time.Second)
	res := p.Check("127.0.0.1", port)

	if !res.Open {
		t.Errorf("expected port %d to be open, got closed (err: %v)", port, res.Err)
	}
	if res.Latency <= 0 {
		t.Error("expected positive latency")
	}
}

func TestProbeClosedPort(t *testing.T) {
	p := NewProbe(200 * time.Millisecond)
	// Port 1 is almost certainly closed and not privileged-accessible in tests.
	res := p.Check("127.0.0.1", 1)
	if res.Open {
		t.Skip("port 1 unexpectedly open, skipping")
	}
	if res.Err == nil {
		t.Error("expected non-nil error for closed port")
	}
}

func TestProbeCheckAll(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port

	p := NewProbe(time.Second)
	results := p.CheckAll("127.0.0.1", []int{port})

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Open {
		t.Errorf("expected open result for port %d", port)
	}
}

func TestProbeDefaultTimeout(t *testing.T) {
	p := NewProbe(0)
	if p.timeout != 2*time.Second {
		t.Errorf("expected default timeout 2s, got %v", p.timeout)
	}
}

func TestProbeSetTimeout(t *testing.T) {
	p := NewProbe(time.Second)
	p.SetTimeout(500 * time.Millisecond)
	if p.timeout != 500*time.Millisecond {
		t.Errorf("expected 500ms, got %v", p.timeout)
	}
	// Non-positive values should be ignored.
	p.SetTimeout(-1)
	if p.timeout != 500*time.Millisecond {
		t.Error("timeout should not change on non-positive value")
	}
}
