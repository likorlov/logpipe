package sink_test

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/your/logpipe"
	"github.com/your/logpipe/sink"
)

type countingSink struct {
	calls   atomic.Int32
	failFor int32 // fail the first N calls
}

func (c *countingSink) Write(_ logpipe.Entry) error {
	n := c.calls.Add(1)
	if n <= c.failFor {
		return errors.New("transient error")
	}
	return nil
}

func (c *countingSink) Close() error { return nil }

func TestRetrySink_SucceedsFirstTry(t *testing.T) {
	cs := &countingSink{}
	s := sink.NewRetrySink(cs, sink.RetryOptions{MaxAttempts: 3, Delay: 0})
	if err := s.Write(logpipe.Entry{Message: "ok"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cs.calls.Load() != 1 {
		t.Fatalf("expected 1 call, got %d", cs.calls.Load())
	}
}

func TestRetrySink_RetriesAndSucceeds(t *testing.T) {
	cs := &countingSink{failFor: 2}
	s := sink.NewRetrySink(cs, sink.RetryOptions{MaxAttempts: 3, Delay: time.Millisecond})
	if err := s.Write(logpipe.Entry{Message: "hi"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cs.calls.Load() != 3 {
		t.Fatalf("expected 3 calls, got %d", cs.calls.Load())
	}
}

func TestRetrySink_ExhaustsAttempts(t *testing.T) {
	cs := &countingSink{failFor: 10}
	s := sink.NewRetrySink(cs, sink.RetryOptions{MaxAttempts: 3, Delay: 0})
	err := s.Write(logpipe.Entry{Message: "fail"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if cs.calls.Load() != 3 {
		t.Fatalf("expected 3 calls, got %d", cs.calls.Load())
	}
}

func TestRetrySink_Backoff(t *testing.T) {
	cs := &countingSink{failFor: 2}
	start := time.Now()
	s := sink.NewRetrySink(cs, sink.RetryOptions{
		MaxAttempts: 3,
		Delay:       10 * time.Millisecond,
		Multiplier:  2.0,
	})
	_ = s.Write(logpipe.Entry{})
	if elapsed := time.Since(start); elapsed < 30*time.Millisecond {
		t.Fatalf("expected backoff delay >= 30ms, got %v", elapsed)
	}
}
