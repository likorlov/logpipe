package sink_test

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/your/logpipe"
	"github.com/your/logpipe/sink"
)

func TestQuotaSink_AllowsUnderQuota(t *testing.T) {
	collect := &captureSink{}
	q := sink.NewQuotaSink(collect, 3, time.Second, nil)

	for i := 0; i < 3; i++ {
		if err := q.Write(logpipe.Entry{Message: "ok"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if collect.count() != 3 {
		t.Fatalf("expected 3 entries, got %d", collect.count())
	}
}

func TestQuotaSink_DropsOverQuota(t *testing.T) {
	collect := &captureSink{}
	q := sink.NewQuotaSink(collect, 2, time.Second, nil)

	for i := 0; i < 5; i++ {
		q.Write(logpipe.Entry{Message: "msg"}) //nolint:errcheck
	}
	if collect.count() != 2 {
		t.Fatalf("expected 2 forwarded, got %d", collect.count())
	}
}

func TestQuotaSink_ResetsAfterWindow(t *testing.T) {
	collect := &captureSink{}
	q := sink.NewQuotaSink(collect, 1, 50*time.Millisecond, nil)

	if err := q.Write(logpipe.Entry{Message: "first"}); err != nil {
		t.Fatal(err)
	}
	if err := q.Write(logpipe.Entry{Message: "over"}); err == nil {
		t.Fatal("expected quota error")
	}
	time.Sleep(60 * time.Millisecond)
	if err := q.Write(logpipe.Entry{Message: "reset"}); err != nil {
		t.Fatalf("expected reset, got: %v", err)
	}
	if collect.count() != 2 {
		t.Fatalf("expected 2 entries, got %d", collect.count())
	}
}

func TestQuotaSink_PerKeyQuota(t *testing.T) {
	collect := &captureSink{}
	keyFn := func(e logpipe.Entry) string { return e.Fields["svc"] }
	q := sink.NewQuotaSink(collect, 1, time.Second, keyFn)

	q.Write(logpipe.Entry{Fields: map[string]any{"svc": "a"}}) //nolint:errcheck
	q.Write(logpipe.Entry{Fields: map[string]any{"svc": "b"}}) //nolint:errcheck
	if collect.count() != 2 {
		t.Fatalf("expected 2 (one per key), got %d", collect.count())
	}
}

func TestQuotaSink_Close(t *testing.T) {
	closed := false
	cs := &callbackSink{closeFn: func() error { closed = true; return nil }}
	q := sink.NewQuotaSink(cs, 10, time.Second, nil)
	if err := q.Close(); err != nil {
		t.Fatal(err)
	}
	if !closed {
		t.Fatal("expected inner sink to be closed")
	}
}

// helpers shared across quota tests

type captureSink struct{ n atomic.Int32 }

func (c *captureSink) Write(logpipe.Entry) error { c.n.Add(1); return nil }
func (c *captureSink) Close() error              { return nil }
func (c *captureSink) count() int                { return int(c.n.Load()) }

type callbackSink struct {
	writeFn func(logpipe.Entry) error
	closeFn func() error
}

func (s *callbackSink) Write(e logpipe.Entry) error {
	if s.writeFn != nil {
		return s.writeFn(e)
	}
	return nil
}
func (s *callbackSink) Close() error {
	if s.closeFn != nil {
		return s.closeFn()
	}
	return nil
}

var _ = errors.New // suppress unused import
