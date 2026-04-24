package sink

import (
	"math/rand"
	"time"

	"github.com/logpipe/logpipe"
)

// JitterSink wraps an inner Sink and introduces a random delay before forwarding
// each log entry. This is useful for spreading bursts of log writes across time
// to avoid thundering-herd effects on downstream sinks.
//
// The delay is uniformly distributed in [0, maxJitter).
type JitterSink struct {
	inner     logpipe.Sink
	maxJitter time.Duration
	rng       *rand.Rand
}

// NewJitterSink returns a JitterSink that delays each Write by a random
// duration in [0, maxJitter) before forwarding to inner.
// Panics if maxJitter <= 0.
func NewJitterSink(inner logpipe.Sink, maxJitter time.Duration) *JitterSink {
	if maxJitter <= 0 {
		panic("logpipe/sink: JitterSink maxJitter must be positive")
	}
	return &JitterSink{
		inner:     inner,
		maxJitter: maxJitter,
		//nolint:gosec // non-cryptographic jitter is intentional
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Write sleeps for a random duration in [0, maxJitter) then forwards entry
// to the inner sink.
func (s *JitterSink) Write(entry logpipe.Entry) error {
	delay := time.Duration(s.rng.Int63n(int64(s.maxJitter)))
	time.Sleep(delay)
	return s.inner.Write(entry)
}

// Close closes the inner sink.
func (s *JitterSink) Close() error {
	return s.inner.Close()
}
