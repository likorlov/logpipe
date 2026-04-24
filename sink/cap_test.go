package sink_test

import (
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestCapSink_AllowsUnderCap(t *testing.T) {
	var received []logpipe.Entry
	inner := &callbackSink{fn: func(e logpipe.Entry) error { received = append(received, e); return nil }}
	s := sink.NewCapSink(inner, 3)

	for i := 0; i < 3; i++ {
		if err := s.Write(logpipe.Entry{"n": i}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if len(received) != 3 {
		t.Fatalf("expected 3 forwarded entries, got %d", len(received))
	}
}

func TestCapSink_DropsOverCap(t *testing.T) {
	var received []logpipe.Entry
	inner := &callbackSink{fn: func(e logpipe.Entry) error { received = append(received, e); return nil }}
	s := sink.NewCapSink(inner, 2)

	for i := 0; i < 5; i++ {
		if err := s.Write(logpipe.Entry{"n": i}); err != nil {
			t.Fatalf("unexpected error on write %d: %v", i, err)
		}
	}
	if len(received) != 2 {
		t.Fatalf("expected 2 forwarded entries, got %d", len(received))
	}
	if s.Count() != 5 {
		t.Fatalf("expected Count 5, got %d", s.Count())
	}
}

func TestCapSink_Reset(t *testing.T) {
	var received []logpipe.Entry
	inner := &callbackSink{fn: func(e logpipe.Entry) error { received = append(received, e); return nil }}
	s := sink.NewCapSink(inner, 1)

	_ = s.Write(logpipe.Entry{"a": 1})
	_ = s.Write(logpipe.Entry{"a": 2}) // dropped

	s.Reset()
	_ = s.Write(logpipe.Entry{"a": 3}) // should be forwarded

	if len(received) != 2 {
		t.Fatalf("expected 2 entries after reset, got %d", len(received))
	}
}

func TestCapSink_PanicOnZeroCap(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for max=0")
		}
	}()
	sink.NewCapSink(&callbackSink{}, 0)
}

func TestCapSink_Close(t *testing.T) {
	closed := false
	inner := &callbackSink{closeFn: func() error { closed = true; return nil }}
	s := sink.NewCapSink(inner, 10)
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
	if !closed {
		t.Fatal("inner sink was not closed")
	}
}
