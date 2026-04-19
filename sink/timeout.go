package sink

import (
	"context"
	"fmt"
	"time"

	"github.com/logpipe/logpipe"
)

// TimeoutSink wraps a sink and enforces a per-write deadline.
// If the inner sink does not return within the timeout, the write
// is abandoned and an error is returned to the caller.
type TimeoutSink struct {
	inner   logpipe.Sink
	timeout time.Duration
}

// NewTimeoutSink returns a sink that cancels writes to inner that
// exceed the given timeout duration.
func NewTimeoutSink(inner logpipe.Sink, timeout time.Duration) *TimeoutSink {
	return &TimeoutSink{inner: inner, timeout: timeout}
}

func (s *TimeoutSink) Write(entry logpipe.Entry) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	type result struct{ err error }
	ch := make(chan result, 1)

	go func() {
		ch <- result{err: s.inner.Write(entry)}
	}()

	select {
	case r := <-ch:
		return r.err
	case <-ctx.Done():
		return fmt.Errorf("logpipe/sink: write timed out after %s", s.timeout)
	}
}

func (s *TimeoutSink) Close() error {
	return s.inner.Close()
}
