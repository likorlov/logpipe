package sink

import (
	"sync"
	"time"

	"github.com/your/logpipe"
)

// RateLimitSink drops log entries that exceed a maximum count per interval.
type RateLimitSink struct {
	mu       sync.Mutex
	wrapped  logpipe.Sink
	max      int
	interval time.Duration
	count    int
	reset    time.Time
}

// NewRateLimitSink returns a Sink that forwards at most max entries per
// interval to wrapped, silently dropping entries that exceed the limit.
func NewRateLimitSink(wrapped logpipe.Sink, max int, interval time.Duration) *RateLimitSink {
	return &RateLimitSink{
		wrapped:  wrapped,
		max:      max,
		interval: interval,
		reset:    time.Now().Add(interval),
	}
}

func (r *RateLimitSink) Write(entry logpipe.Entry) error {
	r.mu.Lock()
	now := time.Now()
	if now.After(r.reset) {
		r.count = 0
		r.reset = now.Add(r.interval)
	}
	if r.count >= r.max {
		r.mu.Unlock()
		return nil
	}
	r.count++
	r.mu.Unlock()
	return r.wrapped.Write(entry)
}

func (r *RateLimitSink) Close() error {
	return r.wrapped.Close()
}
