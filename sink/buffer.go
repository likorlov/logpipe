package sink

import (
	"sync"
	"time"

	"github.com/example/logpipe"
)

// BufferedSink wraps another sink and batches entries, flushing on interval or size.
type BufferedSink struct {
	inner    logpipe.Sink
	buf      []logpipe.Entry
	mu       sync.Mutex
	maxSize  int
	interval time.Duration
	stopCh   chan struct{}
	wg       sync.WaitGroup
}

// NewBufferedSink creates a BufferedSink that flushes when buf hits maxSize or interval elapses.
func NewBufferedSink(inner logpipe.Sink, maxSize int, interval time.Duration) *BufferedSink {
	b := &BufferedSink{
		inner:    inner,
		maxSize:  maxSize,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
	b.wg.Add(1)
	go b.run()
	return b
}

func (b *BufferedSink) Write(e logpipe.Entry) error {
	b.mu.Lock()
	b.buf = append(b.buf, e)
	should := len(b.buf) >= b.maxSize
	b.mu.Unlock()
	if should {
		return b.Flush()
	}
	return nil
}

// Flush writes all buffered entries to the inner sink.
func (b *BufferedSink) Flush() error {
	b.mu.Lock()
	entries := b.buf
	b.buf = nil
	b.mu.Unlock()
	for _, e := range entries {
		if err := b.inner.Write(e); err != nil {
			return err
		}
	}
	return nil
}

func (b *BufferedSink) run() {
	defer b.wg.Done()
	ticker := time.NewTicker(b.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			_ = b.Flush()
		case <-b.stopCh:
			_ = b.Flush()
			return
		}
	}
}

func (b *BufferedSink) Close() error {
	close(b.stopCh)
	b.wg.Wait()
	return b.inner.Close()
}
