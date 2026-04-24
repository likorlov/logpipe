package sink

import (
	"errors"
	"sync"
	"time"

	"github.com/logpipe/logpipe"
)

// ExpireSink drops entries whose timestamp field is older than the configured
// max age. Entries missing the timestamp field are forwarded unchanged.
type expireSink struct {
	inner    logpipe.Sink
	field    string
	maxAge   time.Duration
	nowFn    func() time.Time
	mu       sync.Mutex
}

// NewExpireSink returns a Sink that discards log entries older than maxAge.
// The age is determined by the value of field, which must be a time.Time.
// Use the functional option WithExpireField to override the default field
// name ("ts").
func NewExpireSink(inner logpipe.Sink, maxAge time.Duration, opts ...func(*expireSink)) logpipe.Sink {
	s := &expireSink{
		inner:  inner,
		field:  "ts",
		maxAge: maxAge,
		nowFn:  time.Now,
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// WithExpireField overrides the entry field used to read the timestamp.
func WithExpireField(field string) func(*expireSink) {
	return func(s *expireSink) { s.field = field }
}

func (s *expireSink) Write(entry logpipe.Entry) error {
	s.mu.Lock()
	now := s.nowFn()
	s.mu.Unlock()

	val, ok := entry.Fields[s.field]
	if !ok {
		return s.inner.Write(entry)
	}

	ts, ok := val.(time.Time)
	if !ok {
		return s.inner.Write(entry)
	}

	if now.Sub(ts) > s.maxAge {
		return nil // drop expired entry
	}

	return s.inner.Write(entry)
}

func (s *expireSink) Close() error {
	return s.inner.Close()
}

var _ logpipe.Sink = (*expireSink)(nil)
var errExpired = errors.New("entry expired")
_ = errExpired // kept for future use
