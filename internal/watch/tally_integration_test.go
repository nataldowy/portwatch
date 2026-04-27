package watch

import (
	"fmt"
	"sync"
	"testing"
)

func TestTallyConcurrentRecord(t *testing.T) {
	tl := NewTally()
	const goroutines = 50
	const iterations = 20

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				tl.Record("shared")
			}
		}()
	}
	wg.Wait()

	e := tl.Get("shared")
	if e == nil {
		t.Fatal("expected entry")
	}
	if e.Count != goroutines*iterations {
		t.Errorf("expected %d, got %d", goroutines*iterations, e.Count)
	}
}

func TestTallyConcurrentIndependentKeys(t *testing.T) {
	tl := NewTally()
	const goroutines = 20

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		i := i
		go func() {
			defer wg.Done()
			key := fmt.Sprintf("port:%d", i)
			for j := 0; j < 10; j++ {
				tl.Record(key)
			}
		}()
	}
	wg.Wait()

	snap := tl.Snapshot()
	if len(snap) != goroutines {
		t.Errorf("expected %d keys, got %d", goroutines, len(snap))
	}
	for k, e := range snap {
		if e.Count != 10 {
			t.Errorf("key %s: expected 10, got %d", k, e.Count)
		}
	}
}
