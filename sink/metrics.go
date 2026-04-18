package sink

import (
	"sync/atomic"

	"github.com/andybar2/logpipe"
)

// MetricsSink wraps a Sink and tracks write counts, drop counts, and errors.
type MetricsSink struct {
	next    logpipe.Sink
	writes  atomic.Int64
	drops   atomic.Int64
	errors  atomic.Int64
}

// NewMetricsSink returns a MetricsSink wrapping next.
func NewMetricsSink(next logpipe.Sink) *MetricsSink {
	return &MetricsSink{next: next}
}

// Write forwards the entry to the wrapped sink, recording outcomes.
func (m *MetricsSink) Write(e logpipe.Entry) error {
	err := m.next.Write(e)
	if err != nil {
		if err == ErrDropped {
			m.drops.Add(1)
		} else {
			m.errors.Add(1)
		}
		return err
	}
	m.writes.Add(1)
	return nil
}

// Close closes the wrapped sink.
func (m *MetricsSink) Close() error {
	return m.next.Close()
}

// Writes returns the number of successfully written entries.
func (m *MetricsSink) Writes() int64 { return m.writes.Load() }

// Drops returns the number of dropped entries (ErrDropped).
func (m *MetricsSink) Drops() int64 { return m.drops.Load() }

// Errors returns the number of write errors (excluding drops).
func (m *MetricsSink) Errors() int64 { return m.errors.Load() }

// Reset zeroes all counters.
func (m *MetricsSink) Reset() {
	m.writes.Store(0)
	m.drops.Store(0)
	m.errors.Store(0)
}
