package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestEnrichSink_InjectsFields(t *testing.T) {
	var got logpipe.Entry
	inner := &captureSink{writeFn: func(e logpipe.Entry) error { got = e; return nil }}

	s := sink.NewEnrichSink(inner, func() map[string]any {
		return map[string]any{"host": "srv-1", "version": "1.0"}
	})

	_ = s.Write(logpipe.Entry{Level: logpipe.Info, Message: "hi", Fields: map[string]any{}})

	if got.Fields["host"] != "srv-1" {
		t.Errorf("expected host=srv-1, got %v", got.Fields["host"])
	}
	if got.Fields["version"] != "1.0" {
		t.Errorf("expected version=1.0, got %v", got.Fields["version"])
	}
}

func TestEnrichSink_EntryFieldsWin(t *testing.T) {
	var got logpipe.Entry
	inner := &captureSink{writeFn: func(e logpipe.Entry) error { got = e; return nil }}

	s := sink.NewEnrichSink(inner, func() map[string]any {
		return map[string]any{"host": "srv-1"}
	})

	_ = s.Write(logpipe.Entry{
		Level:   logpipe.Info,
		Message: "hi",
		Fields:  map[string]any{"host": "override"},
	})

	if got.Fields["host"] != "override" {
		t.Errorf("expected host=override, got %v", got.Fields["host"])
	}
}

func TestEnrichSink_NoMutationOfOriginal(t *testing.T) {
	inner := &captureSink{writeFn: func(e logpipe.Entry) error { return nil }}
	s := sink.NewEnrichSink(inner, func() map[string]any {
		return map[string]any{"injected": true}
	})

	orig := logpipe.Entry{Fields: map[string]any{"a": 1}}
	_ = s.Write(orig)

	if _, ok := orig.Fields["injected"]; ok {
		t.Error("original entry was mutated")
	}
}

func TestEnrichSink_PropagatesError(t *testing.T) {
	want := errors.New("sink error")
	inner := &captureSink{writeFn: func(e logpipe.Entry) error { return want }}
	s := sink.NewEnrichSink(inner, func() map[string]any { return nil })

	if err := s.Write(logpipe.Entry{}); !errors.Is(err, want) {
		t.Errorf("expected %v, got %v", want, err)
	}
}

func TestEnrichSink_Close(t *testing.T) {
	closed := false
	inner := &captureSink{closeFn: func() error { closed = true; return nil }}
	s := sink.NewEnrichSink(inner, func() map[string]any { return nil })
	_ = s.Close()
	if !closed {
		t.Error("expected inner sink to be closed")
	}
}
