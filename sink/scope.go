package sink

import (
	"github.com/logpipe/logpipe"
)

// ScopeSink injects a fixed "scope" (or namespace) field into every log entry
// before forwarding it to the inner sink. This is useful for tagging all logs
// produced by a subsystem with a common identifier.
type ScopeSink struct {
	inner logpipe.Sink
	field string
	scope string
}

// NewScopeSink returns a ScopeSink that sets entry.Fields[field] = scope on
// every entry, then forwards to inner. If field is empty, "scope" is used.
func NewScopeSink(inner logpipe.Sink, scope, field string) *ScopeSink {
	if field == "" {
		field = "scope"
	}
	return &ScopeSink{inner: inner, scope: scope, field: field}
}

// Write injects the scope field and forwards the entry.
func (s *ScopeSink) Write(entry logpipe.Entry) error {
	copy := logpipe.Entry{
		Level:   entry.Level,
		Message: entry.Message,
		Fields:  make(map[string]any, len(entry.Fields)+1),
	}
	for k, v := range entry.Fields {
		copy.Fields[k] = v
	}
	if _, exists := copy.Fields[s.field]; !exists {
		copy.Fields[s.field] = s.scope
	}
	return s.inner.Write(copy)
}

// Close closes the inner sink.
func (s *ScopeSink) Close() error {
	return s.inner.Close()
}
