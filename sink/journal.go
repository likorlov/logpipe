package sink

import (
	"fmt"
	"sync"

	"github.com/your-org/logpipe"
)

// JournalSink records every entry written to it along with a monotonically
// increasing journal index, then forwards the enriched entry to an inner sink.
// It is useful for audit trails and ordered replay scenarios.
type JournalSink struct {
	mu    sync.Mutex
	inner logpipe.Sink
	field string
	index uint64
}

// NewJournalSink creates a JournalSink that injects a journal sequence number
// into field (default "_journal") before forwarding to inner.
func NewJournalSink(inner logpipe.Sink, field string) *JournalSink {
	if field == "" {
		field = "_journal"
	}
	return &JournalSink{inner: inner, field: field}
}

// Write stamps the entry with the next journal index and forwards it.
func (j *JournalSink) Write(entry logpipe.Entry) error {
	if entry == nil {
		return fmt.Errorf("journal: nil entry")
	}

	j.mu.Lock()
	j.index++
	idx := j.index
	j.mu.Unlock()

	out := make(logpipe.Entry, len(entry)+1)
	for k, v := range entry {
		out[k] = v
	}
	out[j.field] = idx

	return j.inner.Write(out)
}

// Index returns the current journal index (number of entries written so far).
func (j *JournalSink) Index() uint64 {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.index
}

// Reset resets the journal index back to zero.
func (j *JournalSink) Reset() {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.index = 0
}

// Close closes the inner sink.
func (j *JournalSink) Close() error {
	return j.inner.Close()
}
