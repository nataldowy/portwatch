package watch

import "sync"

// Buffer accumulates events up to a maximum size, dropping oldest when full.
type Buffer struct {
	mu    sync.Mutex
	items []BufferEntry
	max   int
}

// BufferEntry holds a key and associated count.
type BufferEntry struct {
	Key   string
	Count int
}

// NewBuffer creates a Buffer with the given maximum capacity.
func NewBuffer(max int) *Buffer {
	if max <= 0 {
		max = 64
	}
	return &Buffer{max: max}
}

// Add increments the count for key, or appends a new entry.
// If at capacity the oldest entry is dropped.
func (b *Buffer) Add(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i := range b.items {
		if b.items[i].Key == key {
			b.items[i].Count++
			return
		}
	}
	if len(b.items) >= b.max {
		b.items = b.items[1:]
	}
	b.items = append(b.items, BufferEntry{Key: key, Count: 1})
}

// Flush returns all entries and clears the buffer.
func (b *Buffer) Flush() [].mu.Lock()
	defer b.mu.Unlock()
	out := make([]BufferEntry, len(b.items))
	copy(out, b.items)
	b.items = b.items[:0]
	return out
}

// Len returns the current number of distinct keys.
func (b *Buffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.items)
}

// Reset clears the buffer.
func (b *Buffer) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.items = b.items[:0]
}
