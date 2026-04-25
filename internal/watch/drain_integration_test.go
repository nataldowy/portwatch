package watch

import (
	"sync"
	"testing"
	"time"
)

func TestDrainConcurrentAdd(t *testing.T) {
	d := NewDrain(time.Minute, 1000)
	now := time.Now()
	const goroutines = 20
	const addsEach = 25

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < addsEach; j++ {
				d.Add("tcp:80", now)
			}
		}()
	}
	wg.Wait()

	got := d.Count("tcp:80")
	if got != goroutines*addsEach {
		t.Fatalf("expected %d events, got %d", goroutines*addsEach, got)
	}
}

func TestDrainConcurrentFlushAndAdd(t *testing.T) {
	d := NewDrain(time.Minute, 1000)
	now := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d.Add("k", now)
		}()
	}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d.Flush("k", now)
		}()
	}
	wg.Wait()
	// No panic or race — correctness is the goal here.
}
