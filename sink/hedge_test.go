package sink_test

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestHedgeSink_PrimaryFastNoHedge(t *testing.T) {
	var secondaryCalled atomic.Bool
	primary := &callbackSink{fn: func(logpipe.Entry) error { return nil }}
	secondary := &callbackSink{fn: func(logpipe.Entry) error {
		secondaryCalled.Store(true)
		return nil
	}}

	h := sink.NewHedgeSink(primary, secondary, 50*time.Millisecond)
	if err := h.Write(logpipe.Entry{"msg": "hello"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	time.Sleep(80 * time.Millisecond)
	if secondaryCalled.Load() {
		t.Error("secondary should not have been called when primary was fast")
	}
}

func TestHedgeSink_SecondaryFiredOnSlowPrimary(t *testing.T) {
	var secondaryCalled atomic.Bool
	primary := &callbackSink{fn: func(logpipe.Entry) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	}}
	secondary := &callbackSink{fn: func(logpipe.Entry) error {
		secondaryCalled.Store(true)
		return nil
	}}

	h := sink.NewHedgeSink(primary, secondary, 20*time.Millisecond)
	if err := h.Write(logpipe.Entry{"msg": "slow"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !secondaryCalled.Load() {
		t.Error("secondary should have been called after hedge delay")
	}
}

func TestHedgeSink_BothFail(t *testing.T) {
	errP := errors.New("primary error")
	errS := errors.New("secondary error")
	primary := &callbackSink{fn: func(logpipe.Entry) error {
		time.Sleep(30 * time.Millisecond)
		return errP
	}}
	secondary := &callbackSink{fn: func(logpipe.Entry) error { return errS }}

	h := sink.NewHedgeSink(primary, secondary, 10*time.Millisecond)
	err := h.Write(logpipe.Entry{"msg": "fail"})
	if err == nil {
		t.Fatal("expected an error when both sinks fail")
	}
}

func TestHedgeSink_PanicOnZeroDelay(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for zero delay")
		}
	}()
	sink.NewHedgeSink(
		&callbackSink{fn: func(logpipe.Entry) error { return nil }},
		&callbackSink{fn: func(logpipe.Entry) error { return nil }},
		0,
	)
}

func TestHedgeSink_Close(t *testing.T) {
	closed := 0
	mk := func() *callbackSink {
		return &callbackSink{
			fn:      func(logpipe.Entry) error { return nil },
			closeFn: func() error { closed++; return nil },
		}
	}
	h := sink.NewHedgeSink(mk(), mk(), 10*time.Millisecond)
	if err := h.Close(); err != nil {
		t.Fatalf("close error: %v", err)
	}
	if closed != 2 {
		t.Errorf("expected 2 closes, got %d", closed)
	}
}

// callbackSink is a test helper already defined in other _test files in this
// package; if not, define a minimal version here.
type callbackSink struct {
	fn      func(logpipe.Entry) error
	closeFn func() error
}

func (c *callbackSink) Write(e logpipe.Entry) error { return c.fn(e) }
func (c *callbackSink) Close() error {
	if c.closeFn != nil {
		return c.closeFn()
	}
	return nil
}
