package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

type recordSink struct {
	entries []logpipe.Entry
	writeErr error
	closed bool
}

func (r *recordSink) Write(e logpipe.Entry) error {
	if r.writeErr != nil {
		return r.writeErr
	}
	r.entries = append(r.entries, e)
	return nil
}

func (r *recordSink) Close() error { r.closed = true; return nil }

func TestFallbackSink_PrimarySucceeds(t *testing.T) {
	primary := &recordSink{}
	fallback := &recordSink{}
	s := sink.NewFallbackSink(primary, fallback)

	entry := logpipe.Entry{Message: "hello"}
	if err := s.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(primary.entries) != 1 {
		t.Errorf("expected primary to receive entry")
	}
	if len(fallback.entries) != 0 {
		t.Errorf("expected fallback to be unused")
	}
}

func TestFallbackSink_UsedOnPrimaryError(t *testing.T) {
	primary := &recordSink{writeErr: errors.New("primary down")}
	fallback := &recordSink{}
	s := sink.NewFallbackSink(primary, fallback)

	entry := logpipe.Entry{Message: "hello"}
	if err := s.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fallback.entries) != 1 {
		t.Errorf("expected fallback to receive entry")
	}
}

func TestFallbackSink_BothFail(t *testing.T) {
	primary := &recordSink{writeErr: errors.New("primary down")}
	fallback := &recordSink{writeErr: errors.New("fallback down")}
	s := sink.NewFallbackSink(primary, fallback)

	if err := s.Write(logpipe.Entry{Message: "x"}); err == nil {
		t.Error("expected error when both sinks fail")
	}
}

func TestFallbackSink_Close(t *testing.T) {
	primary := &recordSink{}
	fallback := &recordSink{}
	s := sink.NewFallbackSink(primary, fallback)

	if err := s.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
	if !primary.closed || !fallback.closed {
		t.Error("expected both sinks to be closed")
	}
}
