package sink

import (
	"math/rand"
	"sync"

	"github.com/yourorg/logpipe"
)

// SamplingSink wraps a Sink and forwards only a fraction of log entries.
// A rate of 1.0 forwards all entries; 0.0 drops all entries.
type SamplingSink struct {
	mu   sync.Mutex
	inner logpipe.Sink
	rate  float64
	rng   *rand.Rand
}

// NewSamplingSink creates a SamplingSink that forwards entries to inner
// with the given sample rate (0.0–1.0).
func NewSamplingSink(inner logpipe.Sink, rate float64, src rand.Source) *SamplingSink {
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}
	if src == nil {
		src = rand.NewSource(42)
	}
	return &SamplingSink{
		inner: inner,
		rate:  rate,
		rng:   rand.New(src),
	}
}

// Write forwards the entry to the inner sink only if it passes the sample check.
func (s *SamplingSink) Write(entry logpipe.Entry) error {
	s.mu.Lock()
	pass := s.rng.Float64() < s.rate
	s.mu.Unlock()
	if !pass {
		return nil
	}
	return s.inner.Write(entry)
}

// Close closes the inner sink.
func (s *SamplingSink) Close() error {
	return s.inner.Close()
}
