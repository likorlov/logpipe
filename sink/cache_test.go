package sink_test

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestCacheSink_ForwardsFirstEntry(t *testing.T) {
	var calls int32
	inner := &callCountSink{fn: func(e logpipe.Entry) error {
		atomic.AddInt32(&calls, 1)
		return nil
	}}
	s := sink.NewCacheSink(inner, 100*time.Millisecond)
	defer s.Close()

	entry := logpipe.Entry{Fields: logpipe.Fields{"message": "hello"}}
	if err := s.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if atomic.LoadInt32(&calls) != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestCacheSink_DeduplicatesWithinTTL(t *testing.T) {
	var calls int32
	inner := &callCountSink{fn: func(e logpipe.Entry) error {
		atomic.AddInt32(&calls, 1)
		return nil
	}}
	s := sink.NewCacheSink(inner, 200*time.Millisecond)
	defer s.Close()

	entry := logpipe.Entry{Fields: logpipe.Fields{"message": "repeat"}}
	for i := 0; i < 5; i++ {
		if err := s.Write(entry); err != nil {
			t.Fatalf("unexpected error on write %d: %v", i, err)
		}
	}
	if atomic.LoadInt32(&calls) != 1 {
		t.Fatalf("expected 1 inner call, got %d", calls)
	}
}

func TestCacheSink_ForwardsAfterTTLExpires(t *testing.T) {
	var calls int32
	inner := &callCountSink{fn: func(e logpipe.Entry) error {
		atomic.AddInt32(&calls, 1)
		return nil
	}}
	s := sink.NewCacheSink(inner, 20*time.Millisecond)
	defer s.Close()

	entry := logpipe.Entry{Fields: logpipe.Fields{"message": "ttl-test"}}
	_ = s.Write(entry)
	time.Sleep(40 * time.Millisecond)
	_ = s.Write(entry)

	if atomic.LoadInt32(&calls) != 2 {
		t.Fatalf("expected 2 inner calls after TTL, got %d", calls)
	}
}

func TestCacheSink_CachesError(t *testing.T) {
	sentinel := errors.New("inner failure")
	var calls int32
	inner := &callCountSink{fn: func(e logpipe.Entry) error {
		atomic.AddInt32(&calls, 1)
		return sentinel
	}}
	s := sink.NewCacheSink(inner, 200*time.Millisecond)
	defer s.Close()

	entry := logpipe.Entry{Fields: logpipe.Fields{"message": "bad"}}
	err1 := s.Write(entry)
	err2 := s.Write(entry)

	if !errors.Is(err1, sentinel) || !errors.Is(err2, sentinel) {
		t.Fatalf("expected sentinel error, got %v / %v", err1, err2)
	}
	if atomic.LoadInt32(&calls) != 1 {
		t.Fatalf("expected 1 inner call, got %d", calls)
	}
}

func TestCacheSink_NoMessageFieldAlwaysForwards(t *testing.T) {
	var calls int32
	inner := &callCountSink{fn: func(e logpipe.Entry) error {
		atomic.AddInt32(&calls, 1)
		return nil
	}}
	s := sink.NewCacheSink(inner, time.Second)
	defer s.Close()

	entry := logpipe.Entry{Fields: logpipe.Fields{"level": "info"}}
	_ = s.Write(entry)
	_ = s.Write(entry)

	if atomic.LoadInt32(&calls) != 2 {
		t.Fatalf("expected 2 inner calls for keyless entries, got %d", calls)
	}
}

// callCountSink is a test helper that delegates to fn.
type callCountSink struct {
	fn func(logpipe.Entry) error
}

func (c *callCountSink) Write(e logpipe.Entry) error { return c.fn(e) }
func (c *callCountSink) Close() error               { return nil }
