package sink

import (
	"time"

	"github.com/logpipe/logpipe"
)

// TimestampSink injects a timestamp field into every log entry before
// forwarding it to the inner sink. If the entry already contains the
// target field it is overwritten.
type TimestampSink struct {
	inner logpipe.Sink
	field string
	now   func() time.Time
}

// NewTimestampSink returns a sink that stamps each entry with the current
// UTC time under field. Pass an empty field name to use the default "ts".
func NewTimestampSink(inner logpipe.Sink, field string) *TimestampSink {
	if field == "" {
		field = "ts"
	}
	return &TimestampSink{inner: inner, field: field, now: func() time.Time {
		return time.Now().UTC()
	}}
}

func (s *TimestampSink) Write(e logpipe.Entry) error {
	stamped := logpipe.Entry{
		Level:   e.Level,
		Message: e.Message,
		Fields:  make(map[string]any, len(e.Fields)+1),
	}
	for k, v := range e.Fields {
		stamped.Fields[k] = v
	}
	stamped.Fields[s.field] = s.now().Format(time.RFC3339Nano)
	return s.inner.Write(stamped)
}

func (s *TimestampSink) Close() error {
	return s.inner.Close()
}
