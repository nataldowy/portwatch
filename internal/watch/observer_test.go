package watch

import (
	"sync"
	"testing"
	"time"
)

func TestObserverSubscribeAndEmit(t *testing.T) {
	obs := NewObserver(0)
	var got []ObserverEvent
	obs.Subscribe("scan", func(e ObserverEvent) { got = append(got, e) })
	obs.Emit("scan", "payload-1")
	obs.Emit("scan", "payload-2")
	if len(got) != 2 {
		t.Fatalf("expected 2 events, got %d", len(got))
	}
	if got[0].Payload != "payload-1" {
		t.Errorf("unexpected payload: %v", got[0].Payload)
	}
}

func TestObserverNoSubscribersNoError(t *testing.T) {
	obs := NewObserver(0)
	// should not panic
	obs.Emit("unknown", nil)
}

func TestObserverMultipleSubscribers(t *testing.T) {
	obs := NewObserver(0)
	count := 0
	obs.Subscribe("tick", func(ObserverEvent) { count++ })
	obs.Subscribe("tick", func(ObserverEvent) { count++ })
	obs.Emit("tick", nil)
	if count != 2 {
		t.Errorf("expected count=2, got %d", count)
	}
}

func TestObserverEventTimestamp(t *testing.T) {
	obs := NewObserver(0)
	before := time.Now()
	var ts time.Time
	obs.Subscribe("ev", func(e ObserverEvent) { ts = e.Timestamp })
	obs.Emit("ev", nil)
	after := time.Now()
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v not in expected range", ts)
	}
}

func TestObserverReset(t *testing.T) {
	obs := NewObserver(0)
	obs.Subscribe("x", func(ObserverEvent) {})
	obs.Reset()
	if obs.SubscriberCount("x") != 0 {
		t.Error("expected 0 subscribers after reset")
	}
}

func TestObserverConcurrentEmit(t *testing.T) {
	obs := NewObserver(0)
	var mu sync.Mutex
	var results []string
	obs.Subscribe("c", func(e ObserverEvent) {
		mu.Lock()
		results = append(results, e.Payload.(string))
		mu.Unlock()
	})
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			obs.Emit("c", "data")
		}()
	}
	wg.Wait()
	if len(results) != 20 {
		t.Errorf("expected 20 results, got %d", len(results))
	}
}
