package sink

import (
	"sync"
	"time"

	"github.com/logpipe/logpipe"
)

// BatchSink accumulates log entries and flushes them as a slice to a
// user-supplied flush function, either when the batch reaches a maximum
// size or when a ticker fires.
type BatchSink struct {
	mu       sync.Mutex
	buf      []logpipe.Entry
	maxSize  int
	flushFn  func([]logpipe.Entry) error
	ticker   *time.Ticker
	done     chan struct{}
}

// NewBatchSink creates a BatchSink that calls flushFn whenever the batch
// reaches maxSize entries or the flushInterval elapses.
func NewBatchSink(maxSize int, flushInterval time.Duration, flushFn func([]logpipe.Entry) error) *BatchSink {
	s := &BatchSink{
		buf:     make([]logpipe.Entry, 0, maxSize),
		maxSize: maxSize,
		flushFn: flushFn,
		ticker:  time.NewTicker(flushInterval),
		done:    make(chan struct{}),
	}
	go s.loop()
	return s
}

func (s *BatchSink) Write(e logpipe.Entry) error {
	s.mu.Lock()
	s.buf = append(s.buf, e)
	ready := len(s.buf) >= s.maxSize
	s.mu.Unlock()
	if ready {
		return s.flush()
	}
	return nil
}

func (s *BatchSink) flush() error {
	s.mu.Lock()
	if len(s.buf) == 0 {
		s.mu.Unlock()
		return nil
	}
	batch := make([]logpipe.Entry, len(s.buf))
	copy(batch, s.buf)
	s.buf = s.buf[:0]
	s.mu.Unlock()
	return s.flushFn(batch)
}

func (s *BatchSink) loop() {
	for {
		select {
		case <-s.ticker.C:
			_ = s.flush()
		case <-s.done:
			return
		}
	}
}

func (s *BatchSink) Close() error {
	s.ticker.Stop()
	close(s.done)
	return s.flush()
}
