package sink

import (
	"github.com/logpipe/logpipe"
)

// CloneSink creates a deep copy of each log entry before forwarding it to the
// inner sink. This is useful when downstream sinks mutate the entry and you
// want to protect the original from modification.
type cloneSink struct {
	inner logpipe.Sink
}

// NewCloneSink wraps inner so that every entry is cloned before being written.
func NewCloneSink(inner logpipe.Sink) logpipe.Sink {
	return &cloneSink{inner: inner}
}

func (s *cloneSink) Write(entry logpipe.Entry) error {
	return s.inner.Write(cloneEntry(entry))
}

func (s *cloneSink) Close() error {
	return s.inner.Close()
}

// cloneEntry returns a shallow copy of entry with a new Fields map.
func cloneEntry(e logpipe.Entry) logpipe.Entry {
	fields := make(map[string]any, len(e.Fields))
	for k, v := range e.Fields {
		fields[k] = v
	}
	return logpipe.Entry{
		Level:   e.Level,
		Message: e.Message,
		Fields:  fields,
	}
}
