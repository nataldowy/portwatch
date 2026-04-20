package watch

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func TestAggregatorConcurrentAdd(t *testing.T) {
	const workers = 10
	const perWorker = 20

	a := NewAggregator(0)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < perWorker; j++ {
				a.Add(scanner.Port{Number: base*100 + j, Proto: "tcp"})
			}
		}(i)
	}

	wg.Wait()
	evts := a.Flush()
	if len(evts) != workers*perWorker {
		t.Fatalf("expected %d events, got %d", workers*perWorker, len(evts))
	}
}

func TestAggregatorMaxUnderConcurrency(t *testing.T) {
	const cap = 50
	a := NewAggregator(cap)
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				a.Add(scanner.Port{Number: base*10 + j, Proto: "tcp"})
			}
		}(i)
	}

	wg.Wait()
	if a.Len() > cap {
		t.Fatalf("aggregator exceeded max capacity: got %d, want <= %d", a.Len(), cap)
	}
}
