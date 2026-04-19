package sink

import (
	"fmt"

	"github.com/logpipe/logpipe"
)

// LabelSink wraps an inner sink and prepends a fixed label to every log
// entry's message field, making it easy to identify log sources in
// aggregated output.
type LabelSink struct {
	inner logpipe.Sink
	label string
	field string
}

// NewLabelSink returns a sink that prepends label to the value of field in
// every entry before forwarding to inner. If field is empty it defaults to
// "message".
func NewLabelSink(inner logpipe.Sink, label, field string) *LabelSink {
	if field == "" {
		field = "message"
	}
	return &LabelSink{inner: inner, label: label, field: field}
}

func (s *LabelSink) Write(entry logpipe.Entry) error {
	copy := logpipe.Entry{
		Level:  entry.Level,
		Fields: make(map[string]any, len(entry.Fields)),
	}
	for k, v := range entry.Fields {
		copy.Fields[k] = v
	}
	if existing, ok := copy.Fields[s.field]; ok {
		copy.Fields[s.field] = fmt.Sprintf("%s %v", s.label, existing)
	} else {
		copy.Fields[s.field] = s.label
	}
	return s.inner.Write(copy)
}

func (s *LabelSink) Close() error {
	return s.inner.Close()
}
