package sink

import (
	"sync"
	"time"

	"github.com/logpipe/logpipe"
)

// throttleSink drops entries for a key if another entry with the same key
// was written within the cooldown window. Unlike rate limiting (which uses
// a token bucket per interval), throttle enforces a minimum gap between
// successive identical messages.
type throttleSink struct {
	next     logpipe.Sink
	keyFn    func(logpipe.Entry) string
	cooldown time.Duration
	mu       sync.Mutex
	last     map[string]time.Time
}

// NewThrottleSink wraps next and suppresses entries whose key (returned by
// keyFn) was seen within cooldown. A nil keyFn falls back to the entry
// message field.
func NewThrottleSink(next logpipe.Sink, cooldown time.Duration, keyFn func(logpipe.Entry) string) logpipe.Sink {
	if keyFn == nil {
		keyFn = func(e logpipe.Entry) string { return e.Message }
	}
	return &throttleSink{
		next:     next,
		keyFn:    keyFn,
		cooldown: cooldown,
		last:     make(map[string]time.Time),
	}
}

func (t *throttleSink) Write(e logpipe.Entry) error {
	key := t.keyFn(e)
	now := time.Now()

	t.mu.Lock()
	if ts, ok := t.last[key]; ok && now.Sub(ts) < t.cooldown {
		t.mu.Unlock()
		return nil
	}
	t.last[key] = now
	t.mu.Unlock()

	return t.next.Write(e)
}

func (t *throttleSink) Close() error {
	return t.next.Close()
}
