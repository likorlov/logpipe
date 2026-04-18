package sink_test

import (
	"errors"
	"testing"
	"time"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestRedactSink_MasksSensitiveFields(t *testing.T) {
	col := &collectSink{}
	s := sink.NewRedactSink(col, "***", "password", "token")

	e := logpipe.Entry{
		Level:   logpipe.LevelInfo,
		Message: "user login",
		Time:    time.Now(),
		Fields:  map[string]any{"user": "alice", "password": "s3cr3t", "token": "abc123"},
	}
	if err := s.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
	got := col.entries[0].Fields
	if got["password"] != "***" {
		t.Errorf("expected password redacted, got %v", got["password"])
	}
	if got["token"] != "***" {
		t.Errorf("expected token redacted, got %v", got["token"])
	}
	if got["user"] != "alice" {
		t.Errorf("expected user preserved, got %v", got["user"])
	}
}

func TestRedactSink_CustomMask(t *testing.T) {
	col := &collectSink{}
	s := sink.NewRedactSink(col, "[REDACTED]", "secret")

	e := logpipe.Entry{
		Level:   logpipe.LevelInfo,
		Message: "test",
		Time:    time.Now(),
		Fields:  map[string]any{"secret": "topsecret"},
	}
	_ = s.Write(e)
	if col.entries[0].Fields["secret"] != "[REDACTED]" {
		t.Errorf("expected custom mask, got %v", col.entries[0].Fields["secret"])
	}
}

func TestRedactSink_NoMutation(t *testing.T) {
	col := &collectSink{}
	s := sink.NewRedactSink(col, "", "key")

	orig := map[string]any{"key": "value", "other": "data"}
	e := logpipe.Entry{Level: logpipe.LevelDebug, Time: time.Now(), Fields: orig}
	_ = s.Write(e)

	if orig["key"] != "value" {
		t.Error("original entry fields were mutated")
	}
}

func TestRedactSink_PropagatesError(t *testing.T) {
	errSink := &errWriteSink{err: errors.New("write failed")}
	s := sink.NewRedactSink(errSink, "", "x")
	e := logpipe.Entry{Level: logpipe.LevelInfo, Time: time.Now(), Fields: map[string]any{"x": 1}}
	if err := s.Write(e); err == nil {
		t.Error("expected error from wrapped sink")
	}
}

func TestRedactSink_Close(t *testing.T) {
	col := &collectSink{}
	s := sink.NewRedactSink(col, "", "k")
	if err := s.Close(); err != nil {
		t.Errorf("unexpected close error: %v", err)
	}
}
