package sink

import (
	"sync"

	"github.com/your/logpipe"
)

// SnapshotSink keeps the last N log entries in memory for inspection.
// It is useful for testing and diagnostics.
type SnapshotSink struct {
	mu      sync.RWMutex
	entries []logpipe.Entry
	max     int
}

// NewSnapshotSink creates a SnapshotSink that retains up to max entries.
// When the buffer is full, the oldest entry is evicted.
func NewSnapshotSink(max int) *SnapshotSink {
	if max <= 0 {
		max = 100
	}
	return &SnapshotSink{max: max}
}

// Write stores the entry in the ring buffer.
func (s *SnapshotSink) Write(e logpipe.Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.entries) >= s.max {
		s.entries = s.entries[1:]
	}
	s.entries = append(s.entries, e)
	return nil
}

// Entries returns a shallow copy of all retained entries.
func (s *SnapshotSink) Entries() []logpipe.Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]logpipe.Entry, len(s.entries))
	copy(out, s.entries)
	return out
}

// Len returns the number of retained entries.
func (s *SnapshotSink) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}

// Reset clears all retained entries.
func (s *SnapshotSink) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = s.entries[:0]
}

// Close is a no-op.
func (s *SnapshotSink) Close() error { return nil }
