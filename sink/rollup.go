package sink

import (
	"sync"
	"time"

	"github.com/logpipe/logpipe"
)

// rollupSink accumulates numeric field values over a time window and flushes
// a single aggregated entry (with sum, count, min, max) to the inner sink.
type rollupSink struct {
	inner  logpipe.Sink
	field  string
	window time.Duration

	mu    sync.Mutex
	sum   float64
	count int64
	min   float64
	max   float64
	first bool // true when no entries seen yet in current window

	ticker *time.Ticker
	done   chan struct{}
}

// NewRollupSink returns a Sink that accumulates values of the given numeric
// field over window, then forwards one summary entry per window to inner.
// The summary entry contains fields: "sum", "count", "min", "max", and
// "field" set to the monitored field name.
func NewRollupSink(inner logpipe.Sink, field string, window time.Duration) logpipe.Sink {
	if window <= 0 {
		panic("rollup: window must be positive")
	}
	s := &rollupSink{
		inner:  inner,
		field:  field,
		window: window,
		first:  true,
		ticker: time.NewTicker(window),
		done:   make(chan struct{}),
	}
	go s.loop()
	return s
}

func (s *rollupSink) Write(e logpipe.Entry) error {
	v, ok := e.Fields[s.field]
	if !ok {
		return nil
	}
	f, ok := toFloat(v)
	if !ok {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sum += f
	s.count++
	if s.first || f < s.min {
		s.min = f
	}
	if s.first || f > s.max {
		s.max = f
	}
	s.first = false
	return nil
}

func (s *rollupSink) loop() {
	for {
		select {
		case <-s.ticker.C:
			s.flush()
		case <-s.done:
			s.ticker.Stop()
			s.flush()
			return
		}
	}
}

func (s *rollupSink) flush() {
	s.mu.Lock()
	if s.count == 0 {
		s.mu.Unlock()
		return
	}
	entry := logpipe.Entry{
		Fields: map[string]interface{}{
			"field": s.field,
			"sum":   s.sum,
			"count": s.count,
			"min":   s.min,
			"max":   s.max,
		},
	}
	s.sum = 0
	s.count = 0
	s.min = 0
	s.max = 0
	s.first = true
	s.mu.Unlock()
	_ = s.inner.Write(entry)
}

func (s *rollupSink) Close() error {
	close(s.done)
	return s.inner.Close()
}
