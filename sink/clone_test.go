package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestCloneSink_ForwardsEntry(t *testing.T) {
	var got logpipe.Entry
	inner := &captureSink{fn: func(e logpipe.Entry) error { got = e; return nil }}
	s := sink.NewCloneSink(inner)

	e := logpipe.Entry{Level: logpipe.Info, Message: "hello", Fields: map[string]any{"k": "v"}}
	if err := s.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Message != "hello" || got.Fields["k"] != "v" {
		t.Fatalf("entry not forwarded correctly: %+v", got)
	}
}

func TestCloneSink_NoMutationOfOriginal(t *testing.T) {
	inner := &captureSink{fn: func(e logpipe.Entry) error {
		e.Fields["injected"] = "yes"
		return nil
	}}
	s := sink.NewCloneSink(inner)

	orig := logpipe.Entry{Level: logpipe.Info, Message: "msg", Fields: map[string]any{"a": "1"}}
	if err := s.Write(orig); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := orig.Fields["injected"]; ok {
		t.Fatal("original entry was mutated")
	}
}

func TestCloneSink_PropagatesError(t *testing.T) {
	expected := errors.New("sink error")
	inner := &captureSink{fn: func(e logpipe.Entry) error { return expected }}
	s := sink.NewCloneSink(inner)

	err := s.Write(logpipe.Entry{Fields: map[string]any{}})
	if !errors.Is(err, expected) {
		t.Fatalf("expected %v, got %v", expected, err)
	}
}

func TestCloneSink_Close(t *testing.T) {
	closed := false
	inner := &captureSink{fn: func(e logpipe.Entry) error { return nil }, closeFn: func() error { closed = true; return nil }}
	s := sink.NewCloneSink(inner)
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !closed {
		t.Fatal("inner sink was not closed")
	}
}

// captureSink is a test helper that calls fn on each Write.
type captureSink struct {
	fn      func(logpipe.Entry) error
	closeFn func() error
}

func (c *captureSink) Write(e logpipe.Entry) error { return c.fn(e) }
func (c *captureSink) Close() error {
	if c.closeFn != nil {
		return c.closeFn()
	}
	return nil
}
