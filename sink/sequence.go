package sink

import (
	"fmt"
	"sync/atomic"

	"github.com/your/logpipe"
)

// SequenceSink wraps an inner Sink and injects a monotonically increasing
// sequence number into every log entry under a configurable field name.
type SequenceSink struct {
	inner logpipe.Sink
	field string
	counter uint64
}

// NewSequenceSink returns a SequenceSink that stamps each entry with a
// sequence number stored in field (default "seq" when empty).
func NewSequenceSink(inner logpipe.Sink, field string) *SequenceSink {
	if field == "" {
		field = "seq"
	}
	return &SequenceSink{inner: inner, field: field}
}

func (s *SequenceSink) Write(entry logpipe.Entry) error {
	n := atomic.AddUint64(&s.counter, 1)

	// Copy fields to avoid mutating the original entry.
	fields := make(map[string]any, len(entry.Fields)+1)
	for k, v := range entry.Fields {
		fields[k] = v
	}
	fields[s.field] = fmt.Sprintf("%d", n)

	entry.Fields = fields
	return s.inner.Write(entry)
}

func (s *SequenceSink) Close() error { return s.inner.Close() }

// Counter returns the current sequence value.
func (s *SequenceSink) Counter() uint64 { return atomic.LoadUint64(&s.counter) }
