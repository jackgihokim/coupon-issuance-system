package id

import "sync"

type ID struct {
	mu sync.Mutex
	n  uint64
}

// NewID creates and returns a pointer to a new ID instance initialized with zero.
func NewID() *ID {
	return &ID{n: 0}
}

// Next generates and returns the next sequential uint64 value in a thread-safe manner.
func (i *ID) Next() uint64 {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.n++
	return i.n
}
