package sink

import (
	"sync"
	"time"

	"github.com/logpipe/logpipe"
)

// WindowSink collects log entries within a sliding time window and forwards
// only entries that fall within the window duration relative to now.
type WindowSink struct {
	inner    logpipe.Sink
	window   time.Duration
	mu       sync.Mutex
	entries  []timedEntry
}

type timedEntry struct {
	at    time.Time
	entry logpipe.Entry
}

// NewWindowSink returns a sink that only forwards entries received within
// the given sliding window duration. Older entries are evicted on each write.
func NewWindowSink(inner logpipe.Sink, window time.Duration) *WindowSink {
	return &WindowSink{inner: inner, window: window}
}

func (w *WindowSink) Write(entry logpipe.Entry) error {
	now := time.Now()
	cutoff := now.Add(-w.window)

	w.mu.Lock()
	// evict expired entries
	valid := w.entries[:0]
	for _, e := range w.entries {
		if e.at.After(cutoff) {
			valid = append(valid, e)
		}
	}
	w.entries = append(valid, timedEntry{at: now, entry: entry})
	w.mu.Unlock()

	return w.inner.Write(entry)
}

// Entries returns a snapshot of entries currently within the window.
func (w *WindowSink) Entries() []logpipe.Entry {
	now := time.Now()
	cutoff := now.Add(-w.window)

	w.mu.Lock()
	defer w.mu.Unlock()

	var out []logpipe.Entry
	for _, e := range w.entries {
		if e.at.After(cutoff) {
			out = append(out, e.entry)
		}
	}
	return out
}

func (w *WindowSink) Close() error { return w.inner.Close() }
