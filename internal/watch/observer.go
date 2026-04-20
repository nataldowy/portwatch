package watch

import (
	"sync"
	"time"
)

// ObserverEvent represents a named event emitted by the observer.
type ObserverEvent struct {
	Name      string
	Payload   interface{}
	Timestamp time.Time
}

// ObserverHandler is a callback invoked when an event is emitted.
type ObserverHandler func(ObserverEvent)

// Observer is a simple pub/sub event bus for internal watch events.
type Observer struct {
	mu       sync.RWMutex
	subs     map[string][]ObserverHandler
	bufSize  int
}

// NewObserver creates an Observer. bufSize controls the async dispatch
// channel capacity (0 = synchronous).
func NewObserver(bufSize int) *Observer {
	if bufSize < 0 {
		bufSize = 0
	}
	return &Observer{
		subs:    make(map[string][]ObserverHandler),
		bufSize: bufSize,
	}
}

// Subscribe registers a handler for the given event name.
func (o *Observer) Subscribe(name string, h ObserverHandler) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.subs[name] = append(o.subs[name], h)
}

// Emit dispatches an event to all subscribers registered under name.
func (o *Observer) Emit(name string, payload interface{}) {
	ev := ObserverEvent{Name: name, Payload: payload, Timestamp: time.Now()}
	o.mu.RLock()
	handlers := make([]ObserverHandler, len(o.subs[name]))
	copy(handlers, o.subs[name])
	o.mu.RUnlock()

	for _, h := range handlers {
		h(ev)
	}
}

// Reset removes all subscribers.
func (o *Observer) Reset() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.subs = make(map[string][]ObserverHandler)
}

// SubscriberCount returns the number of subscribers for a given event name.
func (o *Observer) SubscriberCount(name string) int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return len(o.subs[name])
}
