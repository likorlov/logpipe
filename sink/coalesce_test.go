package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestCoalesceSink_FirstSucceeds(t *testing.T) {
	var got []logpipe.Entry
	a := &captureSink{fn: func(e logpipe.Entry) error { got = append(got, e); return nil }}
	b := &captureSink{fn: func(e logpipe.Entry) error { t.Fatal("should not reach b"); return nil }}

	s := sink.NewCoalesceSink(a, b)
	entry := logpipe.Entry{Message: "hello", Fields: map[string]any{}}
	if err := s.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 entry in a, got %d", len(got))
	}
}

func TestCoalesceSink_FallsThrough(t *testing.T) {
	failErr := errors.New("fail")
	var got []logpipe.Entry
	a := &captureSink{fn: func(e logpipe.Entry) error { return failErr }}
	b := &captureSink{fn: func(e logpipe.Entry) error { got = append(got, e); return nil }}

	s := sink.NewCoalesceSink(a, b)
	entry := logpipe.Entry{Message: "hi", Fields: map[string]any{}}
	if err := s.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected entry in b, got %d", len(got))
	}
}

func TestCoalesceSink_AllFail(t *testing.T) {
	errA := errors.New("errA")
	errB := errors.New("errB")
	a := &captureSink{fn: func(e logpipe.Entry) error { return errA }}
	b := &captureSink{fn: func(e logpipe.Entry) error { return errB }}

	s := sink.NewCoalesceSink(a, b)
	entry := logpipe.Entry{Message: "oops", Fields: map[string]any{}}
	err := s.Write(entry)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, errA) || !errors.Is(err, errB) {
		t.Fatalf("expected both errors joined, got: %v", err)
	}
}

func TestCoalesceSink_Close(t *testing.T) {
	closed := 0
	mk := func() *captureSink {
		return &captureSink{closeFn: func() error { closed++; return nil }}
	}
	s := sink.NewCoalesceSink(mk(), mk(), mk())
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
	if closed != 3 {
		t.Fatalf("expected 3 closes, got %d", closed)
	}
}

// captureSink is a test helper (local to coalesce tests).
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
