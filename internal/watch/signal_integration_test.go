package watch

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSignalConcurrentWaitersAllUnblock(t *testing.T) {
	s := NewSignal()
	const goroutines = 50
	var wg sync.WaitGroup
	var unblocked int64

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			<-s.Wait()
			atomic.AddInt64(&unblocked, 1)
		}()
	}

	// Give goroutines time to park on Wait.
	time.Sleep(10 * time.Millisecond)
	s.Fire()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("only %d/%d goroutines unblocked", atomic.LoadInt64(&unblocked), goroutines)
	}

	if got := atomic.LoadInt64(&unblocked); got != goroutines {
		t.Fatalf("expected %d unblocked, got %d", goroutines, got)
	}
}

func TestSignalConcurrentFireIsRaceFree(t *testing.T) {
	s := NewSignal()
	var wg sync.WaitGroup
	const goroutines = 20

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			s.Fire()
		}()
	}
	wg.Wait()

	if !s.Fired() {
		t.Fatal("expected signal to be fired after concurrent Fire calls")
	}
}
