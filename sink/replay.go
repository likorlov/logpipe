package sink

import (
	"sync"

	"github.com/logpipe/logpipe"
)

// ReplaySink records every entry it receives and can replay them into any
// other Sink on demand. Useful for testing and deferred processing.
type ReplaySink struct {
	mu      sync.Mutex
	entries []logpipe.Entry
	inner   logpipe.Sink
}

// NewReplaySink returns a ReplaySink that optionally forwards each entry to
// inner (may be nil to record only).
func NewReplaySink(inner logpipe.Sink) *ReplaySink {
	return &ReplaySink{inner: inner}
}

// Write records the entry and, if an inner sink was provided, forwards it.
func (r *ReplaySink) Write(e logpipe.Entry) error {
	r.mu.Lock()
	r.entries = append(r.entries, e)
	r.mu.Unlock()

	if r.inner != nil {
		return r.inner.Write(e)
	}
	return nil
}

// Replay writes all recorded entries into dst in the order they were received.
func (r *ReplaySink) Replay(dst logpipe.Sink) error {
	r.mu.Lock()
	copy := append([]logpipe.Entry(nil), r.entries...)
	r.mu.Unlock()

	for _, e := range copy {
		if err := dst.Write(e); err != nil {
			return err
		}
	}
	return nil
}

// Entries returns a snapshot of all recorded entries.
func (r *ReplaySink) Entries() []logpipe.Entry {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]logpipe.Entry(nil), r.entries...)
}

// Reset clears all recorded entries.
func (r *ReplaySink) Reset() {
	r.mu.Lock()
	r.entries = nil
	r.mu.Unlock()
}

// Close closes the inner sink if one was provided.
func (r *ReplaySink) Close() error {
	if r.inner != nil {
		return r.inner.Close()
	}
	return nil
}
