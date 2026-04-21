package watch

import (
	"net"
	"sync"
	"testing"
	"time"
)

// TestProbeConcurrentCheckAll verifies that CheckAll is safe under concurrent
// calls sharing the same Probe instance.
func TestProbeConcurrentCheckAll(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()
	// Accept connections in the background so the listener doesn't block.
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()

	port := ln.Addr().(*net.TCPAddr).Port
	p := NewProbe(time.Second)

	const goroutines = 10
	var wg sync.WaitGroup
	errCh := make(chan string, goroutines)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results := p.CheckAll("127.0.0.1", []int{port})
			if len(results) != 1 || !results[0].Open {
				errCh <- "expected open port result"
			}
		}()
	}
	wg.Wait()
	close(errCh)

	for msg := range errCh {
		t.Error(msg)
	}
}

// TestProbeSetTimeoutConcurrent ensures SetTimeout is race-free.
func TestProbeSetTimeoutConcurrent(t *testing.T) {
	p := NewProbe(time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			p.SetTimeout(time.Duration(n+1) * 100 * time.Millisecond)
		}(i)
	}
	wg.Wait()
	if p.timeout <= 0 {
		t.Error("timeout should be positive after concurrent sets")
	}
}
