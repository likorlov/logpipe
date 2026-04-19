package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestScopeSink_InjectsScope(t *testing.T) {
	col := &collectSink{}
	s := sink.NewScopeSink(col, "auth", "scope")

	_ = s.Write(logpipe.Entry{Level: logpipe.INFO, Message: "login", Fields: map[string]any{}})

	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
	if col.entries[0].Fields["scope"] != "auth" {
		t.Errorf("expected scope=auth, got %v", col.entries[0].Fields["scope"])
	}
}

func TestScopeSink_DefaultField(t *testing.T) {
	col := &collectSink{}
	s := sink.NewScopeSink(col, "payments", "")

	_ = s.Write(logpipe.Entry{Level: logpipe.INFO, Message: "charge", Fields: map[string]any{}})

	if col.entries[0].Fields["scope"] != "payments" {
		t.Errorf("expected default field 'scope', got %v", col.entries[0].Fields["scope"])
	}
}

func TestScopeSink_EntryFieldOverridesScope(t *testing.T) {
	col := &collectSink{}
	s := sink.NewScopeSink(col, "auth", "scope")

	_ = s.Write(logpipe.Entry{
		Level:   logpipe.INFO,
		Message: "override",
		Fields:  map[string]any{"scope": "custom"},
	})

	if col.entries[0].Fields["scope"] != "custom" {
		t.Errorf("entry field should override sink scope, got %v", col.entries[0].Fields["scope"])
	}
}

func TestScopeSink_NoMutationOfOriginal(t *testing.T) {
	col := &collectSink{}
	s := sink.NewScopeSink(col, "svc", "scope")

	original := logpipe.Entry{Level: logpipe.DEBUG, Message: "test", Fields: map[string]any{}}
	_ = s.Write(original)

	if _, ok := original.Fields["scope"]; ok {
		t.Error("original entry should not be mutated")
	}
}

func TestScopeSink_PropagatesError(t *testing.T) {
	errSink := &errOnWriteSink{err: errors.New("write failed")}
	s := sink.NewScopeSink(errSink, "svc", "scope")

	err := s.Write(logpipe.Entry{Level: logpipe.INFO, Message: "x", Fields: map[string]any{}})
	if err == nil {
		t.Error("expected error to propagate")
	}
}

func TestScopeSink_Close(t *testing.T) {
	col := &collectSink{}
	s := sink.NewScopeSink(col, "svc", "scope")
	if err := s.Close(); err != nil {
		t.Errorf("unexpected close error: %v", err)
	}
}
