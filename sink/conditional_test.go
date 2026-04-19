package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestConditionalSink_RoutesToThen(t *testing.T) {
	then := &captureSink{}
	else_ := &captureSink{}
	cs := sink.NewConditionalSink(func(e logpipe.Entry) bool {
		return e.Fields["branch"] == "then"
	}, then, else_)

	cs.Write(logpipe.Entry{Fields: map[string]any{"branch": "then"}})
	cs.Write(logpipe.Entry{Fields: map[string]any{"branch": "else"}})

	if len(then.entries) != 1 {
		t.Fatalf("expected 1 then entry, got %d", len(then.entries))
	}
	if len(else_.entries) != 1 {
		t.Fatalf("expected 1 else entry, got %d", len(else_.entries))
	}
}

func TestConditionalSink_NilThenDrops(t *testing.T) {
	else_ := &captureSink{}
	cs := sink.NewConditionalSink(func(e logpipe.Entry) bool { return true }, nil, else_)
	if err := cs.Write(logpipe.Entry{Fields: map[string]any{}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(else_.entries) != 0 {
		t.Fatal("else sink should not have received entry")
	}
}

func TestConditionalSink_NilElseDrops(t *testing.T) {
	then := &captureSink{}
	cs := sink.NewConditionalSink(func(e logpipe.Entry) bool { return false }, then, nil)
	if err := cs.Write(logpipe.Entry{Fields: map[string]any{}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(then.entries) != 0 {
		t.Fatal("then sink should not have received entry")
	}
}

func TestConditionalSink_PropagatesError(t *testing.T) {
	want := errors.New("boom")
	then := &errorSink{err: want}
	cs := sink.NewConditionalSink(func(e logpipe.Entry) bool { return true }, then, nil)
	if err := cs.Write(logpipe.Entry{Fields: map[string]any{}}); !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestConditionalSink_Close(t *testing.T) {
	then := &captureSink{}
	else_ := &captureSink{}
	cs := sink.NewConditionalSink(func(e logpipe.Entry) bool { return true }, then, else_)
	if err := cs.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
	if !then.closed || !else_.closed {
		t.Fatal("expected both sinks to be closed")
	}
}
