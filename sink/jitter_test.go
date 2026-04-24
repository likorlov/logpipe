package sink_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestJitterSink_ForwardsEntry(t *testing.T) {
	var received []logpipe.Entry
	var mu sync.Mutex
	inner := &callbackSink{
		writeFn: func(e logpipe.Entry) error {
			mu.Lock()
			defer mu.Unlock()
			received = append(received, e)
			return nil
		},
	}

	s := sink.NewJitterSink(inner, 5*time.Millisecond)
	entry := logpipe.Entry{"msg": "hello"}

	if err := s.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(received))
	}
	if received[0]["msg"] != "hello" {
		t.Errorf("unexpected entry: %v", received[0])
	}
}

func TestJitterSink_DelayIsNonNegative(t *testing.T) {
	inner := &callbackSink{writeFn: func(e logpipe.Entry) error { return nil }}
	s := sink.NewJitterSink(inner, 20*time.Millisecond)

	start := time.Now()
	_ = s.Write(logpipe.Entry{"x": "1"})
	elapsed := time.Since(start)

	if elapsed < 0 {
		t.Error("elapsed time should not be negative")
	}
	if elapsed >= 20*time.Millisecond*10 {
		t.Errorf("delay suspiciously large: %v", elapsed)
	}
}

func TestJitterSink_PropagatesError(t *testing.T) {
	sentinel := errors.New("inner error")
	inner := &callbackSink{writeFn: func(e logpipe.Entry) error { return sentinel }}

	s := sink.NewJitterSink(inner, 5*time.Millisecond)
	err := s.Write(logpipe.Entry{"msg": "test"})
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestJitterSink_PanicOnZeroJitter(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for zero maxJitter")
		}
	}()
	sink.NewJitterSink(&callbackSink{}, 0)
}

func TestJitterSink_Close(t *testing.T) {
	closed := false
	inner := &callbackSink{closeFn: func() error { closed = true; return nil }}
	s := sink.NewJitterSink(inner, 1*time.Millisecond)
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !closed {
		t.Error("expected inner sink to be closed")
	}
}
