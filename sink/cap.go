package sink

import (
	"errors"
	"sync/atomic"

	"github.com/logpipe/logpipe"
)

// CapSink limits the total number of log entries forwarded to the inner sink
// over the lifetime of the sink. Once the cap is reached every subsequent
// Write is silently dropped. Close always delegates to the inner sink.
//
// Use Reset to clear the counter and start accepting entries again.
type CapSink struct {
	inner logpipe.Sink
	max   int64
	count atomic.Int64
}

// NewCapSink returns a CapSink that forwards at most max entries to inner.
// Panics if max is less than 1.
func NewCapSink(inner logpipe.Sink, max int64) *CapSink {
	if max < 1 {
		panic("sink: NewCapSink max must be >= 1")
	}
	return &CapSink{inner: inner, max: max}
}

// Write forwards the entry to the inner sink only when the cap has not yet
// been reached. Entries beyond the cap are silently dropped.
func (s *CapSink) Write(entry logpipe.Entry) error {
	n := s.count.Add(1)
	if n > s.max {
		return nil
	}
	return s.inner.Write(entry)
}

// Count returns the number of Write calls made so far (including dropped ones).
func (s *CapSink) Count() int64 {
	return s.count.Load()
}

// Reset zeroes the internal counter so the sink starts accepting entries again.
func (s *CapSink) Reset() {
	s.count.Store(0)
}

// Close closes the inner sink.
func (s *CapSink) Close() error {
	return s.inner.Close()
}

// ErrCapExceeded is returned by Write when the cap has been reached and the
// caller wants an explicit signal rather than a silent drop. Use CapSink
// directly for the silent-drop behaviour.
var ErrCapExceeded = errors.New("sink: cap exceeded")
