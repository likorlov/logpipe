package sink

import (
	"strings"

	"github.com/logpipe/logpipe"
)

// NormalizeSink transforms field keys in each log entry according to a
// normalization function before forwarding to the inner sink. A common use
// case is converting all keys to snake_case or lowercase.
type normalizeSink struct {
	inner  logpipe.Sink
	normFn func(string) string
}

// NewNormalizeSink returns a Sink that applies normFn to every field key in
// the entry before passing it to inner. If normFn is nil, keys are lowercased.
func NewNormalizeSink(inner logpipe.Sink, normFn func(string) string) logpipe.Sink {
	if normFn == nil {
		normFn = strings.ToLower
	}
	return &normalizeSink{inner: inner, normFn: normFn}
}

func (s *normalizeSink) Write(entry logpipe.Entry) error {
	normalized := logpipe.Entry{
		Level:   entry.Level,
		Message: entry.Message,
		Fields:  make(map[string]any, len(entry.Fields)),
	}
	for k, v := range entry.Fields {
		normalized.Fields[s.normFn(k)] = v
	}
	return s.inner.Write(normalized)
}

func (s *normalizeSink) Close() error {
	return s.inner.Close()
}
