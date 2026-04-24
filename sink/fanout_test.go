package sink_test

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

// concurrentSink records writes with a small artificial delay.
type concurrentSink struct {
	delay   time.Duration
	writes  atomic.Int64
	returnErr error
}

func (c *concurrentSink) Write(_ logpipe.Entry) error {
	time.Sleep(c.delay)
	c.writes.Add(1)
	return c.returnErr
}
func (c *concurrentSink) Close() error { return nil }

func TestFanoutSink_WritesAll(t *testing.T) {
	a := &concurrentSink{delay: 10 * time.Millisecond}
	b := &concurrentSink{delay: 10 * time.Millisecond}
	c := &concurrentSink{delay: 10 * time.Millisecond}

	fs := sink.NewFanoutSink(a, b, c)
	entry := logpipe.Entry{"msg": "hello"}

	start := time.Now()
	if err := fs.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	elapsed := time.Since(start)

	// All three sinks ran concurrently so total time should be well under 3×delay.
	if elapsed > 25*time.Millisecond {
		t.Errorf("writes did not appear concurrent: elapsed %v", elapsed)
	}
	for i, s := range []*concurrentSink{a, b, c} {
		if s.writes.Load() != 1 {
			t.Errorf("sink[%d] expected 1 write, got %d", i, s.writes.Load())
		}
	}
}

func TestFanoutSink_CollectsErrors(t *testing.T) {
	boom := errors.New("boom")
	a := &concurrentSink{}
	b := &concurrentSink{returnErr: boom}
	c := &concurrentSink{returnErr: boom}

	fs := sink.NewFanoutSink(a, b, c)
	err := fs.Write(logpipe.Entry{"msg": "x"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if a.writes.Load() != 1 {
		t.Errorf("successful sink should still have been written")
	}
}

func TestFanoutSink_NoSinks(t *testing.T) {
	fs := sink.NewFanoutSink()
	if err := fs.Write(logpipe.Entry{"msg": "x"}); err != nil {
		t.Fatalf("unexpected error with no sinks: %v", err)
	}
}

func TestFanoutSink_Close(t *testing.T) {
	a := &concurrentSink{}
	b := &concurrentSink{}
	fs := sink.NewFanoutSink(a, b)
	if err := fs.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
}
