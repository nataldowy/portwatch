package watch

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func TestRelayConcurrentSubscribeAndForward(t *testing.T) {
	r := NewRelay(nil)
	var total int64

	const workers = 20
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			r.Subscribe(func(ev scanner.DiffEvent) {
				atomic.AddInt64(&total, 1)
			})
		}()
	}
	wg.Wait()

	r.Forward(scanner.DiffEvent{
		Kind: "new",
		Port: scanner.Port{Number: "9090", Protocol: "tcp"},
		At:   time.Now(),
	})

	if int(atomic.LoadInt64(&total)) != workers {
		t.Errorf("expected %d calls, got %d", workers, total)
	}
}

func TestRelayConcurrentForwardIsRaceFree(t *testing.T) {
	r := NewRelay(nil)
	var count int64
	r.Subscribe(func(scanner.DiffEvent) { atomic.AddInt64(&count, 1) })

	const senders = 50
	var wg sync.WaitGroup
	wg.Add(senders)
	for i := 0; i < senders; i++ {
		go func() {
			defer wg.Done()
			r.Forward(scanner.DiffEvent{
				Kind: "new",
				Port: scanner.Port{Number: "80", Protocol: "tcp"},
				At:   time.Now(),
			})
		}()
	}
	wg.Wait()

	if atomic.LoadInt64(&count) != senders {
		t.Errorf("expected %d forwards, got %d", senders, count)
	}
}
