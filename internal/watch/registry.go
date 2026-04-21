package watch

import (
	"fmt"
	"sync"
)

// Registry tracks named components (notifiers, watchers, etc.) by key,
// allowing dynamic registration and lookup at runtime.
type Registry struct {
	mu    sync.RWMutex
	items map[string]any
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		items: make(map[string]any),
	}
}

// Register stores value under key. Returns an error if the key is already taken.
func (r *Registry) Register(key string, value any) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[key]; exists {
		return fmt.Errorf("registry: key %q already registered", key)
	}
	r.items[key] = value
	return nil
}

// Lookup retrieves the value stored under key.
// The second return value is false when the key is absent.
func (r *Registry) Lookup(key string) (any, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	v, ok := r.items[key]
	return v, ok
}

// Unregister removes the entry for key. It is a no-op when the key is absent.
func (r *Registry) Unregister(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.items, key)
}

// Keys returns a snapshot of all registered keys in no particular order.
func (r *Registry) Keys() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	keys := make([]string, 0, len(r.items))
	for k := range r.items {
		keys = append(keys, k)
	}
	return keys
}

// Len returns the number of registered entries.
func (r *Registry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.items)
}

// Reset removes all entries.
func (r *Registry) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.items = make(map[string]any)
}
