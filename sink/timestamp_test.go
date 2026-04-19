package sink_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestTimestampSink_InjectsField(t *testing.T) {
	col := &collectSink{}
	s := sink.NewTimestampSink(col, "ts")

	e := logpipe.Entry{Level: logpipe.Info, Message: "hello", Fields: map[string]any{}}
	if err := s.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
	ts, ok := col.entries[0].Fields["ts"]
	if !ok {
		t.Fatal("expected 'ts' field to be present")
	}
	if _, err := time.Parse(time.RFC3339Nano, ts.(string)); err != nil {
		t.Fatalf("ts field is not a valid RFC3339Nano time: %v", err)
	}
}

func TestTimestampSink_DefaultField(t *testing.T) {
	col := &collectSink{}
	s := sink.NewTimestampSink(col, "")

	_ = s.Write(logpipe.Entry{Level: logpipe.Info, Message: "x", Fields: map[string]any{}})
	if _, ok := col.entries[0].Fields["ts"]; !ok {
		t.Fatal("expected default field 'ts'")
	}
}

func TestTimestampSink_OverwritesExisting(t *testing.T) {
	col := &collectSink{}
	s := sink.NewTimestampSink(col, "ts")

	e := logpipe.Entry{
		Level:   logpipe.Info,
		Message: "msg",
		Fields:  map[string]any{"ts": "old-value"},
	}
	_ = s.Write(e)
	if col.entries[0].Fields["ts"] == "old-value" {
		t.Fatal("expected ts to be overwritten")
	}
}

func TestTimestampSink_NoMutationOfOriginal(t *testing.T) {
	col := &collectSink{}
	s := sink.NewTimestampSink(col, "ts")

	orig := logpipe.Entry{
		Level:   logpipe.Info,
		Message: "msg",
		Fields:  map[string]any{"key": "val"},
	}
	_ = s.Write(orig)
	if _, ok := orig.Fields["ts"]; ok {
		t.Fatal("original entry should not be mutated")
	}
}

func TestTimestampSink_Close(t *testing.T) {
	col := &collectSink{}
	s := sink.NewTimestampSink(col, "ts")
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected error on Close: %v", err)
	}
	if !col.closed {
		t.Fatal("expected inner sink to be closed")
	}
}
