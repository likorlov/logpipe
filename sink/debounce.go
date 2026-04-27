package sink

import (
	"sync"
	"time"

	"github.com/logpipe/logpipe"
)

// debounceSink delays forwarding an entry until no new entries have arrived
// for the configured quiet period. Only the most recent entry is forwarded.
type debounceSink struct {
	inner  logpipe.Sink
	wait   time.Duration
	mu     sync.Mutex
	pending *logpipe.Entry
	timer  *time.Timer
	closed bool
	wg     sync.WaitGroup
}

// NewDebounceSink returns a Sink that suppresses bursts of entries and
// forwards only the last entry received within each quiet period of wait.
// If wait is zero it panics.
func NewDebounceSink(inner logpipe.Sink, wait time.Duration) logpipe.Sink {
	if wait <= 0 {
		panic("logpipe/sink: NewDebounceSink wait must be positive")
	}
	return &debounceSink{inner: inner, wait: wait}
}

func (d *debounceSink) Write(e logpipe.Entry) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed {
		return nil
	}
	copy := cloneEntry(e)
	d.pending = &copy
	if d.timer != nil {
		d.timer.Reset(d.wait)
		return nil
	}
	d.wg.Add(1)
	d.timer = time.AfterFunc(d.wait, func() {
		defer d.wg.Done()
		d.mu.Lock()
		entry := d.pending
		d.pending = nil
		d.timer = nil
		d.mu.Unlock()
		if entry != nil {
			_ = d.inner.Write(*entry)
		}
	})
	return nil
}

func (d *debounceSink) Close() error {
	d.mu.Lock()
	d.closed = true
	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
	entry := d.pending
	d.pending = nil
	d.mu.Unlock()
	d.wg.Wait()
	if entry != nil {
		_ = d.inner.Write(*entry)
	}
	return d.inner.Close()
}
