package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestMergeSink_InjectsFields(t *testing.T) {
	col := &collectSink{}
	s := sink.NewMergeSink(col, map[string]any{"env": "prod", "region": "us-east"})

	_ = s.Write(logpipe.Entry{Level: logpipe.Info, Message: "hello", Fields: map[string]any{"user": "alice"}})

	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
	e := col.entries[0]
	if e.Fields["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", e.Fields["env"])
	}
	if e.Fields["region"] != "us-east" {
		t.Errorf("expected region=us-east, got %v", e.Fields["region"])
	}
	if e.Fields["user"] != "alice" {
		t.Errorf("expected user=alice, got %v", e.Fields["user"])
	}
}

func TestMergeSink_EntryFieldWins(t *testing.T) {
	col := &collectSink{}
	s := sink.NewMergeSink(col, map[string]any{"env": "prod"})

	_ = s.Write(logpipe.Entry{Level: logpipe.Info, Message: "hi", Fields: map[string]any{"env": "staging"}})

	if col.entries[0].Fields["env"] != "staging" {
		t.Errorf("expected entry field to win, got %v", col.entries[0].Fields["env"])
	}
}

func TestMergeSink_DeepMergesNestedMaps(t *testing.T) {
	col := &collectSink{}
	s := sink.NewMergeSink(col, map[string]any{
		"meta": map[string]any{"service": "api", "version": "1"},
	})

	_ = s.Write(logpipe.Entry{
		Level:   logpipe.Info,
		Message: "deep",
		Fields:  map[string]any{"meta": map[string]any{"host": "web-1"}},
	})

	meta, ok := col.entries[0].Fields["meta"].(map[string]any)
	if !ok {
		t.Fatal("expected meta to be map[string]any")
	}
	if meta["service"] != "api" {
		t.Errorf("expected service=api from base, got %v", meta["service"])
	}
	if meta["host"] != "web-1" {
		t.Errorf("expected host=web-1 from override, got %v", meta["host"])
	}
}

func TestMergeSink_NoMutationOfOriginal(t *testing.T) {
	col := &collectSink{}
	s := sink.NewMergeSink(col, map[string]any{"injected": true})

	orig := map[string]any{"msg": "test"}
	_ = s.Write(logpipe.Entry{Level: logpipe.Info, Message: "x", Fields: orig})

	if _, ok := orig["injected"]; ok {
		t.Error("original entry fields were mutated")
	}
}

func TestMergeSink_PropagatesError(t *testing.T) {
	expected := errors.New("write failed")
	s := sink.NewMergeSink(&errSink{err: expected}, map[string]any{"k": "v"})

	if err := s.Write(logpipe.Entry{Level: logpipe.Info, Message: "x"}); !errors.Is(err, expected) {
		t.Errorf("expected propagated error, got %v", err)
	}
}

func TestMergeSink_Close(t *testing.T) {
	col := &collectSink{}
	s := sink.NewMergeSink(col, map[string]any{})
	if err := s.Close(); err != nil {
		t.Errorf("unexpected close error: %v", err)
	}
}

func TestMergeSink_PanicOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil inner sink")
		}
	}()
	sink.NewMergeSink(nil, map[string]any{})
}
