package sink

import (
	"fmt"

	"github.com/logpipe/logpipe"
)

// PrefixSink prepends a static string prefix to a specified field value before
// forwarding the entry to the inner sink. If the field is absent or not a
// string, the entry is forwarded unchanged.
type PrefixSink struct {
	inner  logpipe.Sink
	field  string
	prefix string
}

// NewPrefixSink returns a PrefixSink that prepends prefix to field in every
// log entry before passing it to inner.
//
// Example:
//
//	sink.NewPrefixSink(console, "message", "[APP] ")
func NewPrefixSink(inner logpipe.Sink, field, prefix string) *PrefixSink {
	if field == "" {
		field = "message"
	}
	return &PrefixSink{inner: inner, field: field, prefix: prefix}
}

// Write prepends the configured prefix to the target field and forwards the
// modified entry to the inner sink.
func (s *PrefixSink) Write(entry logpipe.Entry) error {
	out := make(logpipe.Entry, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	if val, ok := out[s.field]; ok {
		if str, ok := val.(string); ok {
			out[s.field] = fmt.Sprintf("%s%s", s.prefix, str)
		}
	}
	return s.inner.Write(out)
}

// Close closes the inner sink.
func (s *PrefixSink) Close() error {
	return s.inner.Close()
}
