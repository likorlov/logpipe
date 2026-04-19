package sink

import (
	"sync/atomic"

	"github.com/logpipe/logpipe"
)

// CounterSink wraps an inner Sink and exposes a running count of entries
// written, dropped (when inner returns an error), and total calls.
type CounterSink struct {
	inner   logpipe.Sink
	writes  atomic.Int64
	drops   atomic.Int64
}

// NewCounterSink returns a CounterSink wrapping inner.
func NewCounterSink(inner logpipe.Sink) *CounterSink {
	return &CounterSink{inner: inner}
}

// Write forwards the entry to the inner sink, incrementing the appropriate counter.
func (c *CounterSink) Write(entry logpipe.Entry) error {
	err := c.inner.Write(entry)
	if err != nil {
		c.drops.Add(1)
	} else {
		c.writes.Add(1)
	}
	return err
}

// Writes returns the number of successfully forwarded entries.
func (c *CounterSink) Writes() int64 { return c.writes.Load() }

// Drops returns the number of entries that resulted in an error from the inner sink.
func (c *CounterSink) Drops() int64 { return c.drops.Load() }

// Total returns Writes + Drops.
func (c *CounterSink) Total() int64 { return c.writes.Load() + c.drops.Load() }

// Reset zeroes all counters.
func (c *CounterSink) Reset() {
	c.writes.Store(0)
	c.drops.Store(0)
}

// Close closes the inner sink.
func (c *CounterSink) Close() error { return c.inner.Close() }
