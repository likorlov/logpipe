package sink

import (
	"github.com/logpipe/logpipe"
)

// FallbackSink tries a primary sink and, on error, forwards the entry to a
// secondary (fallback) sink. This is useful for ensuring log delivery even
// when the primary destination is temporarily unavailable.
type fallbackSink struct {
	primary  logpipe.Sink
	fallback logpipe.Sink
}

// NewFallbackSink returns a Sink that writes to primary and, if primary
// returns an error, writes to fallback instead.
func NewFallbackSink(primary, fallback logpipe.Sink) logpipe.Sink {
	return &fallbackSink{
		primary:  primary,
		fallback: fallback,
	}
}

func (s *fallbackSink) Write(entry logpipe.Entry) error {
	if err := s.primary.Write(entry); err != nil {
		return s.fallback.Write(entry)
	}
	return nil
}

func (s *fallbackSink) Close() error {
	var firstErr error
	if err := s.primary.Close(); err != nil {
		firstErr = err
	}
	if err := s.fallback.Close(); err != nil && firstErr == nil {
		firstErr = err
	}
	return firstErr
}
