package sink

import (
	"sync"

	"github.com/logpipe/logpipe"
)

// AggregateSink collects log entries and periodically flushes them to an inner
// sink as a single merged entry whose fields are combined from all buffered entries.
type AggregateSink struct {
	mu    sync.Mutex
	inner logpipe.Sink
	buf   []logpipe.Entry
	size  int
	field string
}

// NewAggregateSink returns a sink that accumulates up to size entries, then
// merges all buffered entries into one and forwards it to inner. The merged
// entry carries the last entry's level/message; every field from every entry
// is present (later entries win on key collision). The count of merged entries
// is stored under field (default "agg_count").
func NewAggregateSink(inner logpipe.Sink, size int, field string) *AggregateSink {
	if field == "" {
		field = "agg_count"
	}
	return &AggregateSink{inner: inner, size: size, field: field}
}

// Write buffers e and flushes when the buffer reaches capacity.
func (a *AggregateSink) Write(e logpipe.Entry) error {
	a.mu.Lock()
	a.buf = append(a.buf, e)
	ready := len(a.buf) >= a.size
	a.mu.Unlock()
	if ready {
		return a.Flush()
	}
	return nil
}

// Flush merges all buffered entries and forwards the result to the inner sink.
func (a *AggregateSink) Flush() error {
	a.mu.Lock()
	if len(a.buf) == 0 {
		a.mu.Unlock()
		return nil
	}
	entries := a.buf
	a.buf = nil
	a.mu.Unlock()

	merged := logpipe.Entry{
		Level:   entries[len(entries)-1].Level,
		Message: entries[len(entries)-1].Message,
		Fields:  make(map[string]interface{}, len(entries)*4),
	}
	for _, en := range entries {
		for k, v := range en.Fields {
			merged.Fields[k] = v
		}
	}
	merged.Fields[a.field] = len(entries)
	return a.inner.Write(merged)
}

// Close flushes remaining entries and closes the inner sink.
func (a *AggregateSink) Close() error {
	if err := a.Flush(); err != nil {
		return err
	}
	return a.inner.Close()
}
