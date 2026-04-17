package sink

import (
	"sync"
	"time"

	"github.com/your/logpipe"
)

// DedupeSink suppresses duplicate log entries within a configurable window.
// Two entries are considered duplicates if they share the same level and message.
type DedupeSink struct {
	wrapped  logpipe.Sink
	window   time.Duration
	mu       sync.Mutex
	seen     map[string]time.Time
}

// NewDedupeSink wraps the given sink and drops entries whose (level, message)
// pair was already forwarded within the deduplication window.
func NewDedupeSink(wrapped logpipe.Sink, window time.Duration) *DedupeSink {
	return &DedupeSink{
		wrapped: wrapped,
		window:  window,
		seen:    make(map[string]time.Time),
	}
}

func (d *DedupeSink) Write(e logpipe.Entry) error {
	key := e.Level.String() + "\x00" + e.Message

	d.mu.Lock()
	if t, ok := d.seen[key]; ok && time.Since(t) < d.window {
		d.mu.Unlock()
		return nil
	}
	d.seen[key] = time.Now()
	d.mu.Unlock()

	return d.wrapped.Write(e)
}

func (d *DedupeSink) Close() error {
	return d.wrapped.Close()
}
