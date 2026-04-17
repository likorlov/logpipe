package sink

import (
	"fmt"
	"time"

	"github.com/your/logpipe"
)

// RetryOptions configures the retry behaviour.
type RetryOptions struct {
	MaxAttempts int
	Delay       time.Duration
	Multiplier  float64 // backoff multiplier; 1.0 = constant delay
}

type retrySink struct {
	wrapped  logpipe.Sink
	opts     RetryOptions
}

// NewRetrySink wraps a Sink and retries failed Write calls according to opts.
// MaxAttempts must be >= 1. A Multiplier of 0 is treated as 1.0 (no backoff).
func NewRetrySink(s logpipe.Sink, opts RetryOptions) logpipe.Sink {
	if opts.MaxAttempts < 1 {
		opts.MaxAttempts = 1
	}
	if opts.Multiplier <= 0 {
		opts.Multiplier = 1.0
	}
	return &retrySink{wrapped: s, opts: opts}
}

func (r *retrySink) Write(e logpipe.Entry) error {
	delay := r.opts.Delay
	var lastErr error
	for i := 0; i < r.opts.MaxAttempts; i++ {
		if err := r.wrapped.Write(e); err == nil {
			return nil
		} else {
			lastErr = err
		}
		if i < r.opts.MaxAttempts-1 && delay > 0 {
			time.Sleep(delay)
			delay = time.Duration(float64(delay) * r.opts.Multiplier)
		}
	}
	return fmt.Errorf("retrySink: all %d attempts failed: %w", r.opts.MaxAttempts, lastErr)
}

func (r *retrySink) Close() error {
	return r.wrapped.Close()
}
