package sink

import (
	"fmt"

	"github.com/logpipe/logpipe"
)

// HeaderSink injects a fixed string header into a named field of every log
// entry before forwarding it to the wrapped sink. It is useful for stamping
// entries with a service name, environment, or version string.
type HeaderSink struct {
	inner logpipe.Sink
	field string
	value string
}

// NewHeaderSink returns a HeaderSink that sets field to value on every entry
// before passing it to inner. If field is empty, "header" is used.
func NewHeaderSink(inner logpipe.Sink, field, value string) *HeaderSink {
	if field == "" {
		field = "header"
	}
	return &HeaderSink{inner: inner, field: field, value: value}
}

func (s *HeaderSink) Write(e logpipe.Entry) error {
	fields := make(map[string]any, len(e.Fields)+1)
	for k, v := range e.Fields {
		fields[k] = v
	}
	if _, exists := fields[s.field]; !exists {
		fields[s.field] = s.value
	}
	out := logpipe.Entry{
		Level:   e.Level,
		Message: e.Message,
		Fields:  fields,
	}
	return s.inner.Write(out)
}

func (s *HeaderSink) Close() error {
	return s.inner.Close()
}

// String returns a human-readable description of the sink.
func (s *HeaderSink) String() string {
	return fmt.Sprintf("HeaderSink(field=%q value=%q)", s.field, s.value)
}
