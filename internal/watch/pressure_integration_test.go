package watch

import (
	"sync"
	"testing"
	"time"
)

func TestPressureConcurrentRecord(t *testing.T) {
	const goroutines = 20
	const recordsEach = 5

	p := NewPressure(goroutines*recordsEach, time.Minute)

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < recordsEach; j++ {
				p.Record()
			}
		}()
	}
	wg.Wait()

	if got := p.Count(); got != goroutines*recordsEach {
		t.Fatalf("expected %d events, got %d", goroutines*recordsEach, got)
	}
	if !p.High() {
		t.Fatal("expected high after all records")
	}
}

func TestPressureConcurrentHighAndReset(t *testing.T) {
	p := NewPressure(5, time.Minute)

	var wg sync.WaitGroup
	// Writers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				p.Record()
				time.Sleep(time.Millisecond)
			}
		}()
	}
	// Readers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				_ = p.High()
				_ = p.Count()
				time.Sleep(time.Millisecond)
			}
		}()
	}
	// Resetter
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(5 * time.Millisecond)
		p.Reset()
	}()

	wg.Wait() // must not race or deadlock
}
