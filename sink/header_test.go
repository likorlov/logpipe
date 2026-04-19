package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestHeaderSink_InjectsHeader(t *testing.T) {
	col := &collectSink{}
	s := sink.NewHeaderSink(col, "service", "my-svc")

	e := logpipe.Entry{Level: logpipe.Info, Message: "hello", Fields: map[string]any{}}
	if err := s.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
	if got := col.entries[0].Fields["service"]; got != "my-svc" {
		t.Errorf("expected service=my-svc, got %v", got)
	}
}

func TestHeaderSink_DefaultField(t *testing.T) {
	col := &collectSink{}
	s := sink.NewHeaderSink(col, "", "val")

	e := logpipe.Entry{Level: logpipe.Info, Message: "hi", Fields: map[string]any{}}
	_ = s.Write(e)

	if _, ok := col.entries[0].Fields["header"]; !ok {
		t.Error("expected default field 'header' to be set")
	}
}

func TestHeaderSink_EntryFieldWins(t *testing.T) {
	col := &collectSink{}
	s := sink.NewHeaderSink(col, "env", "prod")

	e := logpipe.Entry{
		Level:   logpipe.Info,
		Message: "msg",
		Fields:  map[string]any{"env": "staging"},
	}
	_ = s.Write(e)

	if got := col.entries[0].Fields["env"]; got != "staging" {
		t.Errorf("expected entry field to win, got %v", got)
	}
}

func TestHeaderSink_NoMutationOfOriginal(t *testing.T) {
	col := &collectSink{}
	s := sink.NewHeaderSink(col, "svc", "x")

	orig := map[string]any{"key": "val"}
	e := logpipe.Entry{Level: logpipe.Info, Message: "m", Fields: orig}
	_ = s.Write(e)

	if _, ok := orig["svc"]; ok {
		t.Error("original entry fields were mutated")
	}
}

func TestHeaderSink_PropagatesError(t *testing.T) {
	sentinel := errors.New("sink error")
	s := sink.NewHeaderSink(&errorSink{err: sentinel}, "k", "v")

	err := s.Write(logpipe.Entry{Fields: map[string]any{}})
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestHeaderSink_Close(t *testing.T) {
	col := &collectSink{}
	s := sink.NewHeaderSink(col, "k", "v")
	if err := s.Close(); err != nil {
		t.Errorf("unexpected close error: %v", err)
	}
}
