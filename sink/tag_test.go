package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestTagSink_InjectsTags(t *testing.T) {
	col := &collectSink{}
	s := sink.NewTagSink(col, map[string]any{"env": "prod", "service": "api"})

	_ = s.Write(logpipe.Entry{Message: "hello", Fields: map[string]any{"user": "alice"}})

	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
	e := col.entries[0]
	if e.Fields["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", e.Fields["env"])
	}
	if e.Fields["service"] != "api" {
		t.Errorf("expected service=api, got %v", e.Fields["service"])
	}
	if e.Fields["user"] != "alice" {
		t.Errorf("expected user=alice, got %v", e.Fields["user"])
	}
}

func TestTagSink_EntryFieldsOverrideTags(t *testing.T) {
	col := &collectSink{}
	s := sink.NewTagSink(col, map[string]any{"env": "prod"})

	_ = s.Write(logpipe.Entry{Message: "hi", Fields: map[string]any{"env": "staging"}})

	e := col.entries[0]
	if e.Fields["env"] != "staging" {
		t.Errorf("entry field should override tag: got %v", e.Fields["env"])
	}
}

func TestTagSink_NoMutationOfOriginal(t *testing.T) {
	col := &collectSink{}
	s := sink.NewTagSink(col, map[string]any{"env": "prod"})

	orig := logpipe.Entry{Message: "x", Fields: map[string]any{"a": 1}}
	_ = s.Write(orig)

	if _, ok := orig.Fields["env"]; ok {
		t.Error("original entry fields should not be mutated")
	}
}

func TestTagSink_PropagatesError(t *testing.T) {
	errSink := &errWriteSink{err: errors.New("boom")}
	s := sink.NewTagSink(errSink, map[string]any{"k": "v"})

	err := s.Write(logpipe.Entry{Message: "x"})
	if err == nil {
		t.Error("expected error to be propagated")
	}
}

func TestTagSink_Close(t *testing.T) {
	col := &collectSink{}
	s := sink.NewTagSink(col, nil)
	if err := s.Close(); err != nil {
		t.Errorf("unexpected close error: %v", err)
	}
}
