package sink_test

import (
	"errors"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/logpipe"
	"github.com/yourorg/logpipe/sink"
)

type countSink struct {
	count int64
	err   error
}

func (c *countSink) Write(_ logpipe.Entry) error {
	atomic.AddInt64(&c.count, 1)
	return c.err
}
func (c *countSink) Close() error { return nil }

func TestSamplingSink_AllPass(t *testing.T) {
	inner := &countSink{}
	s := sink.NewSamplingSink(inner, 1.0, rand.NewSource(1))
	for i := 0; i < 100; i++ {
		_ = s.Write(logpipe.Entry{Message: "x"})
	}
	if inner.count != 100 {
		t.Fatalf("expected 100 entries, got %d", inner.count)
	}
}

func TestSamplingSink_NonePass(t *testing.T) {
	inner := &countSink{}
	s := sink.NewSamplingSink(inner, 0.0, rand.NewSource(1))
	for i := 0; i < 100; i++ {
		_ = s.Write(logpipe.Entry{Message: "x"})
	}
	if inner.count != 0 {
		t.Fatalf("expected 0 entries, got %d", inner.count)
	}
}

func TestSamplingSink_PartialRate(t *testing.T) {
	inner := &countSink{}
	s := sink.NewSamplingSink(inner, 0.5, rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10000; i++ {
		_ = s.Write(logpipe.Entry{Message: "x"})
	}
	if inner.count < 4000 || inner.count > 6000 {
		t.Fatalf("expected ~5000 entries, got %d", inner.count)
	}
}

func TestSamplingSink_PropagatesError(t *testing.T) {
	expected := errors.New("sink error")
	inner := &countSink{err: expected}
	s := sink.NewSamplingSink(inner, 1.0, rand.NewSource(1))
	err := s.Write(logpipe.Entry{Message: "x"})
	if !errors.Is(err, expected) {
		t.Fatalf("expected propagated error, got %v", err)
	}
}

func TestSamplingSink_Close(t *testing.T) {
	inner := &countSink{}
	s := sink.NewSamplingSink(inner, 1.0, nil)
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
}
