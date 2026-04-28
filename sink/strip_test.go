package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestStripSink_RemovesFields(t *testing.T) {
	var got logpipe.Entry
	collect := &captureSink{fn: func(e logpipe.Entry) error { got = e; return nil }}

	s := sink.NewStripSink(collect, "password", "token")
	err := s.Write(logpipe.Entry{"msg": "hello", "password": "secret", "token": "abc123", "user": "alice"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got["password"]; ok {
		t.Error("expected 'password' to be stripped")
	}
	if _, ok := got["token"]; ok {
		t.Error("expected 'token' to be stripped")
	}
	if got["user"] != "alice" {
		t.Errorf("expected 'user' to be preserved, got %v", got["user"])
	}
	if got["msg"] != "hello" {
		t.Errorf("expected 'msg' to be preserved, got %v", got["msg"])
	}
}

func TestStripSink_NoMutationOfOriginal(t *testing.T) {
	collect := &captureSink{fn: func(e logpipe.Entry) error { return nil }}
	s := sink.NewStripSink(collect, "secret")

	orig := logpipe.Entry{"msg": "hi", "secret": "shh"}
	_ = s.Write(orig)

	if _, ok := orig["secret"]; !ok {
		t.Error("original entry was mutated")
	}
}

func TestStripSink_NoFieldsReturnsInner(t *testing.T) {
	collect := &captureSink{fn: func(e logpipe.Entry) error { return nil }}
	s := sink.NewStripSink(collect)
	// When no fields are specified, the inner sink is returned directly.
	if s != collect {
		t.Error("expected inner sink to be returned unchanged when no fields given")
	}
}

func TestStripSink_PropagatesError(t *testing.T) {
	want := errors.New("downstream failure")
	errSink := &captureSink{fn: func(e logpipe.Entry) error { return want }}
	s := sink.NewStripSink(errSink, "x")

	if got := s.Write(logpipe.Entry{"x": 1}); !errors.Is(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestStripSink_Close(t *testing.T) {
	closed := false
	inner := &captureSink{
		fn:      func(e logpipe.Entry) error { return nil },
		closeFn: func() error { closed = true; return nil },
	}
	s := sink.NewStripSink(inner, "field")
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
	if !closed {
		t.Error("expected inner sink to be closed")
	}
}

// captureSink is a minimal test helper used across sink tests.
type captureSink struct {
	fn      func(logpipe.Entry) error
	closeFn func() error
}

func (c *captureSink) Write(e logpipe.Entry) error {
	if c.fn != nil {
		return c.fn(e)
	}
	return nil
}

func (c *captureSink) Close() error {
	if c.closeFn != nil {
		return c.closeFn()
	}
	return nil
}
