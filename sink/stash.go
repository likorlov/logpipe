package sink

import (
	"errors"
	"sync"

	"github.com/logpipe/logpipe"
)

// StashSink holds log entries under named keys so they can be retrieved later.
// Entries stored under the same key overwrite the previous value. It is safe
// for concurrent use.
type StashSink struct {
	mu    sync.RWMutex
	store map[string]logpipe.Entry
	keyFn func(logpipe.Entry) string
	inner logpipe.Sink
}

// NewStashSink creates a StashSink that extracts the stash key from each entry
// using keyFn and forwards the entry to inner after storing it. If keyFn
// returns an empty string the entry is forwarded but not stashed.
func NewStashSink(inner logpipe.Sink, keyFn func(logpipe.Entry) string) *StashSink {
	if keyFn == nil {
		panic("logpipe/sink: NewStashSink keyFn must not be nil")
	}
	return &StashSink{
		store: make(map[string]logpipe.Entry),
		keyFn: keyFn,
		inner: inner,
	}
}

// Write stores the entry under the key returned by keyFn, then forwards it.
func (s *StashSink) Write(e logpipe.Entry) error {
	key := s.keyFn(e)
	if key != "" {
		s.mu.Lock()
		s.store[key] = e
		s.mu.Unlock()
	}
	return s.inner.Write(e)
}

// Get returns the most recently stashed entry for key and true, or a zero
// Entry and false if no entry has been stashed under that key.
func (s *StashSink) Get(key string) (logpipe.Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.store[key]
	return e, ok
}

// Delete removes the stashed entry for key.
func (s *StashSink) Delete(key string) {
	s.mu.Lock()
	delete(s.store, key)
	s.mu.Unlock()
}

// Len returns the number of currently stashed entries.
func (s *StashSink) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.store)
}

// Close flushes the stash and closes the inner sink.
func (s *StashSink) Close() error {
	s.mu.Lock()
	s.store = make(map[string]logpipe.Entry)
	s.mu.Unlock()
	if c, ok := s.inner.(interface{ Close() error }); ok {
		return c.Close()
	}
	return nil
}

var _ logpipe.Sink = (*StashSink)(nil)

// ErrKeyNotFound is returned by helper utilities when a stash key is absent.
var ErrKeyNotFound = errors.New("logpipe/sink: stash key not found")
