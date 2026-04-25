package sink

import (
	"fmt"

	"github.com/logpipe/logpipe"
)

// MergeSink merges a set of fixed extra fields into every log entry before
// forwarding to the inner sink. Fields already present in the entry take
// precedence over the merge fields (entry wins).
//
// Unlike EnrichSink (which also injects fields), MergeSink performs a
// deep-merge when a field value is itself a map[string]any, recursively
// combining nested maps rather than replacing them wholesale.
type mergeSink struct {
	inner  logpipe.Sink
	fields map[string]any
}

// NewMergeSink returns a Sink that deep-merges fields into every entry
// before forwarding to inner.
func NewMergeSink(inner logpipe.Sink, fields map[string]any) logpipe.Sink {
	if inner == nil {
		panic("logpipe/sink: NewMergeSink: inner sink must not be nil")
	}
	copy := make(map[string]any, len(fields))
	for k, v := range fields {
		copy[k] = v
	}
	return &mergeSink{inner: inner, fields: copy}
}

func (s *mergeSink) Write(entry logpipe.Entry) error {
	merged := deepMerge(s.fields, entry.Fields)
	return s.inner.Write(logpipe.Entry{Level: entry.Level, Message: entry.Message, Fields: merged})
}

func (s *mergeSink) Close() error {
	return s.inner.Close()
}

// deepMerge returns a new map that is the result of merging base into
// override. Keys in override win; when both values are map[string]any the
// merge recurses.
func deepMerge(base, override map[string]any) map[string]any {
	out := make(map[string]any, len(base)+len(override))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range override {
		if bv, ok := out[k]; ok {
			bMap, bOk := bv.(map[string]any)
			oMap, oOk := v.(map[string]any)
			if bOk && oOk {
				out[k] = deepMerge(bMap, oMap)
				continue
			}
		}
		out[k] = v
	}
	return out
}

var _ fmt.Stringer = (*mergeSink)(nil)

func (s *mergeSink) String() string {
	return fmt.Sprintf("MergeSink(%d fields)", len(s.fields))
}
