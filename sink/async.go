// Package sink provides output sinks for logpipe.
package sink

import (
	"fmt"
	"sync"

	"github.com/logpipe/logpipe"
)

// AsyncSink wraps another Sink and writes log entries asynchronously
// using a background goroutine and an internal channel queue.
type AsyncSink struct {
	inner   logpipe.Sink
	queue   chan logpipe.Entry
	wg      sync.WaitGroup
	once    sync.Once
	ErrFunc func(err error)
}

// NewAsyncSink creates an AsyncSink that dispatches entries to inner
// asynchronously. queueSize controls the depth of the internal channel.
func NewAsyncSink(inner logpipe.Sink, queueSize int) *AsyncSink {
	if queueSize <= 0 {
		queueSize = 64
	}
	as := &AsyncSink{
		inner: inner,
		queue: make(chan logpipe.Entry, queueSize),
	}
	as.wg.Add(1)
	go as.run()
	return as
}

func (a *AsyncSink) run() {
	defer a.wg.Done()
	for entry := range a.queue {
		if err := a.inner.Write(entry); err != nil && a.ErrFunc != nil {
			a.ErrFunc(err)
		}
	}
}

// Write enqueues entry for asynchronous delivery. Returns an error if
// the queue is full.
func (a *AsyncSink) Write(entry logpipe.Entry) error {
	select {
	case a.queue <- entry:
		return nil
	default:
		return fmt.Errorf("async sink queue full (capacity %d)", cap(a.queue))
	}
}

// Close drains the queue and closes the underlying sink.
func (a *AsyncSink) Close() error {
	var err error
	a.once.Do(func() {
		close(a.queue)
		a.wg.Wait()
		err = a.inner.Close()
	})
	return err
}
