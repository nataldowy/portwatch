package watch

import (
	"sync"
	"testing"
)

func TestObserverConcurrentSubscribeAndEmit(t *testing.T) {
	obs := NewObserver(0)
	var mu sync.Mutex
	var received []ObserverEvent

	// Pre-register one subscriber before concurrent activity.
	obs.Subscribe("port", func(e ObserverEvent) {
		mu.Lock()
		received = append(received, e)
		mu.Unlock()
	})

	var wg sync.WaitGroup
	const emitters = 10
	for i := 0; i < emitters; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			obs.Emit("port", n)
		}(i)
	}
	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if len(received) != emitters {
		t.Errorf("expected %d events, got %d", emitters, len(received))
	}
}

func TestObserverSubscriberCountAccurate(t *testing.T) {
	obs := NewObserver(0)
	if obs.SubscriberCount("e") != 0 {
		t.Error("expected 0 initially")
	}
	obs.Subscribe("e", func(ObserverEvent) {})
	obs.Subscribe("e", func(ObserverEvent) {})
	obs.Subscribe("e", func(ObserverEvent) {})
	if obs.SubscriberCount("e") != 3 {
		t.Errorf("expected 3, got %d", obs.SubscriberCount("e"))
	}
	obs.Reset()
	if obs.SubscriberCount("e") != 0 {
		t.Error("expected 0 after reset")
	}
}
