package sink

import (
	"github.com/logpipe/logpipe"
)

// EnrichFunc is a function that returns additional fields to merge into a log entry.
type EnrichFunc func() map[string]any

type enrichSink struct {
	inner   logpipe.Sink
	enrichFn EnrichFunc
}

// NewEnrichSink returns a Sink that merges fields from enrichFn into every
// log entry before forwarding it to inner. Fields already present on the
// entry take precedence over enriched fields.
//
// Example use-cases: injecting a hostname, build version, or request-id
// sourced from a context at construction time.
func NewEnrichSink(inner logpipe.Sink, fn EnrichFunc) logpipe.Sink {
	return &enrichSink{inner: inner, enrichFn: fn}
}

func (s *enrichSink) Write(entry logpipe.Entry) error {
	extra := s.enrichFn()
	merged := make(map[string]any, len(extra)+len(entry.Fields))
	for k, v := range extra {
		merged[k] = v
	}
	// Entry fields win.
	for k, v := range entry.Fields {
		merged[k] = v
	}
	entry.Fields = merged
	return s.inner.Write(entry)
}

func (s *enrichSink) Close() error {
	return s.inner.Close()
}
