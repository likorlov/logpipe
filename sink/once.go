package sink

import (
	"sync"

	"github.com/logpipe/logpipe"
)

// OnceSink forwards only the first log entry that matches the predicate
// (or every entry if predicate is nil) and drops all subsequent entries.
// Useful for capturing first-occurrence errors or startup events.
type OnceSink struct {
	inner logpipe.Sink
	pred  func(logpipe.Entry) bool
	fired bool
	mu    sync.Mutex
}

// NewOnceSink wraps inner and forwards only the first entry for which pred
// returns true. If pred is nil, the first entry is always forwarded.
func NewOnceSink(inner logpipe.Sink, pred func(logpipe.Entry) bool) *OnceSink {
	return &OnceSink{inner: inner, pred: pred}
}

func (s *OnceSink) Write(e logpipe.Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.fired {
		return nil
	}
	if s.pred != nil && !s.pred(e) {
		return nil
	}
	s.fired = true
	return s.inner.Write(e)
}

// Reset allows the sink to fire once more on the next matching entry.
func (s *OnceSink) Reset() {
	s.mu.Lock()
	s.fired = false
	s.mu.Unlock()
}

func (s *OnceSink) Close() error {
	return s.inner.Close()
}
