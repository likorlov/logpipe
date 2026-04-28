package sink

import "github.com/logpipe/logpipe"

// StripSink removes specified fields from each log entry before forwarding
// it to the inner sink. This is useful for removing sensitive, redundant,
// or overly verbose fields from entries in a pipeline.
//
// The original entry is never mutated; a shallow copy is made with the
// specified fields omitted.
type stripSink struct {
	inner  logpipe.Sink
	fields map[string]struct{}
}

// NewStripSink returns a Sink that removes the given field keys from every
// entry before forwarding to inner.
func NewStripSink(inner logpipe.Sink, fields ...string) logpipe.Sink {
	if len(fields) == 0 {
		return inner
	}
	set := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		set[f] = struct{}{}
	}
	return &stripSink{inner: inner, fields: set}
}

func (s *stripSink) Write(entry logpipe.Entry) error {
	copy := make(logpipe.Entry, len(entry))
	for k, v := range entry {
		if _, remove := s.fields[k]; !remove {
			copy[k] = v
		}
	}
	return s.inner.Write(copy)
}

func (s *stripSink) Close() error {
	return s.inner.Close()
}
