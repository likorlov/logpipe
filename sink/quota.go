package sink

import (
	"fmt"
	"sync"
	"time"

	"github.com/your/logpipe"
)

// QuotaSink enforces a maximum number of log entries per key per time window.
// Entries exceeding the quota are dropped until the window resets.
type QuotaSink struct {
	wrapped  logpipe.Sink
	max      int
	window   time.Duration
	keyFn    func(logpipe.Entry) string
	mu       sync.Mutex
	counters map[string]*quotaBucket
}

type quotaBucket struct {
	count  int
	reset  time.Time
}

// NewQuotaSink creates a QuotaSink that allows at most max entries per keyFn
// value within the given window duration. If keyFn is nil, all entries share
// a single global quota.
func NewQuotaSink(wrapped logpipe.Sink, max int, window time.Duration, keyFn func(logpipe.Entry) string) *QuotaSink {
	if keyFn == nil {
		keyFn = func(logpipe.Entry) string { return "__global__" }
	}
	return &QuotaSink{
		wrapped:  wrapped,
		max:      max,
		window:   window,
		keyFn:    keyFn,
		counters: make(map[string]*quotaBucket),
	}
}

func (q *QuotaSink) Write(e logpipe.Entry) error {
	key := q.keyFn(e)
	now := time.Now()

	q.mu.Lock()
	b, ok := q.counters[key]
	if !ok || now.After(b.reset) {
		b = &quotaBucket{reset: now.Add(q.window)}
		q.counters[key] = b
	}
	b.count++
	count := b.count
	q.mu.Unlock()

	if count > q.max {
		return fmt.Errorf("quota exceeded for key %q", key)
	}
	return q.wrapped.Write(e)
}

func (q *QuotaSink) Close() error {
	return q.wrapped.Close()
}
