package sink

import (
	"sync"

	"github.com/your-org/logpipe"
)

// CorrelateSink injects a correlation ID field into every log entry.
// A new correlation ID is generated per "session" and can be rotated
// by calling Rotate. This is useful for tracing a chain of related
// log entries across a request or job lifecycle.
type CorrelateSink struct {
	inner   logpipe.Sink
	field   string
	idFn    func() string
	mu      sync.RWMutex
	current string
}

// NewCorrelateSink wraps inner, injecting the current correlation ID
// into field (default "correlation_id") on every Write call.
// idFn is called once at construction and on each Rotate to produce
// a new ID; if nil, a simple incrementing counter is used.
func NewCorrelateSink(inner logpipe.Sink, field string, idFn func() string) *CorrelateSink {
	if field == "" {
		field = "correlation_id"
	}
	if idFn == nil {
		idFn = defaultCorrelationID()
	}
	s := &CorrelateSink{
		inner: inner,
		field: field,
		idFn:  idFn,
	}
	s.current = s.idFn()
	return s
}

// Write injects the current correlation ID then forwards to the inner sink.
func (s *CorrelateSink) Write(e logpipe.Entry) error {
	s.mu.RLock()
	id := s.current
	s.mu.RUnlock()

	out := make(logpipe.Entry, len(e)+1)
	for k, v := range e {
		out[k] = v
	}
	if _, exists := out[s.field]; !exists {
		out[s.field] = id
	}
	return s.inner.Write(out)
}

// Rotate generates a new correlation ID. Subsequent writes will use
// the new ID. Safe for concurrent use.
func (s *CorrelateSink) Rotate() {
	next := s.idFn()
	s.mu.Lock()
	s.current = next
	s.mu.Unlock()
}

// Current returns the active correlation ID.
func (s *CorrelateSink) Current() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.current
}

// Close closes the inner sink.
func (s *CorrelateSink) Close() error {
	return s.inner.Close()
}

func defaultCorrelationID() func() string {
	var mu sync.Mutex
	var n uint64
	return func() string {
		mu.Lock()
		n++
		v := n
		mu.Unlock()
		return fmt.Sprintf("corr-%d", v)
	}
}
