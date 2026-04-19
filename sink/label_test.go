package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestLabelSink_PrependsLabel(t *testing.T) {
	col := &collectSink{}
	s := sink.NewLabelSink(col, "[svc]", "message")

	_ = s.Write(logpipe.Entry{Level: logpipe.Info, Fields: map[string]any{"message": "hello"}})

	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
	got := col.entries[0].Fields["message"]
	if got != "[svc] hello" {
		t.Errorf("expected '[svc] hello', got %q", got)
	}
}

func TestLabelSink_MissingField(t *testing.T) {
	col := &collectSink{}
	s := sink.NewLabelSink(col, "[svc]", "message")

	_ = s.Write(logpipe.Entry{Level: logpipe.Info, Fields: map[string]any{}})

	got := col.entries[0].Fields["message"]
	if got != "[svc]" {
		t.Errorf("expected '[svc]', got %q", got)
	}
}

func TestLabelSink_DefaultField(t *testing.T) {
	col := &collectSink{}
	s := sink.NewLabelSink(col, "pfx", "")

	_ = s.Write(logpipe.Entry{Level: logpipe.Info, Fields: map[string]any{"message": "hi"}})

	if col.entries[0].Fields["message"] != "pfx hi" {
		t.Error("default field should be 'message'")
	}
}

func TestLabelSink_NoMutationOfOriginal(t *testing.T) {
	col := &collectSink{}
	s := sink.NewLabelSink(col, "X", "message")

	orig := logpipe.Entry{Level: logpipe.Info, Fields: map[string]any{"message": "original"}}
	_ = s.Write(orig)

	if orig.Fields["message"] != "original" {
		t.Error("original entry should not be mutated")
	}
}

func TestLabelSink_PropagatesError(t *testing.T) {
	errSink := &errorSink{err: errors.New("fail")}
	s := sink.NewLabelSink(errSink, "L", "message")

	err := s.Write(logpipe.Entry{Level: logpipe.Info, Fields: map[string]any{}})
	if err == nil {
		t.Error("expected error from inner sink")
	}
}

func TestLabelSink_Close(t *testing.T) {
	col := &collectSink{}
	s := sink.NewLabelSink(col, "L", "message")
	if err := s.Close(); err != nil {
		t.Errorf("unexpected close error: %v", err)
	}
}
