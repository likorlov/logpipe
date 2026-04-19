package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestPrefixSink_PrependsToField(t *testing.T) {
	col := &collectingSink{}
	s := sink.NewPrefixSink(col, "message", "[APP] ")

	_ = s.Write(logpipe.Entry{"message": "hello world"})

	if got := col.entries[0]["message"]; got != "[APP] hello world" {
		t.Fatalf("expected '[APP] hello world', got %q", got)
	}
}

func TestPrefixSink_DefaultField(t *testing.T) {
	col := &collectingSink{}
	s := sink.NewPrefixSink(col, "", ">> ")

	_ = s.Write(logpipe.Entry{"message": "test"})

	if got := col.entries[0]["message"]; got != ">> test" {
		t.Fatalf("expected '>> test', got %q", got)
	}
}

func TestPrefixSink_NonStringFieldUnchanged(t *testing.T) {
	col := &collectingSink{}
	s := sink.NewPrefixSink(col, "count", "pre_")

	_ = s.Write(logpipe.Entry{"count": 42})

	if got := col.entries[0]["count"]; got != 42 {
		t.Fatalf("expected 42, got %v", got)
	}
}

func TestPrefixSink_NoMutationOfOriginal(t *testing.T) {
	col := &collectingSink{}
	s := sink.NewPrefixSink(col, "message", "X-")

	orig := logpipe.Entry{"message": "original"}
	_ = s.Write(orig)

	if orig["message"] != "original" {
		t.Fatal("original entry was mutated")
	}
}

func TestPrefixSink_PropagatesError(t *testing.T) {
	errSink := &errorSink{err: errors.New("write failed")}
	s := sink.NewPrefixSink(errSink, "message", "pre_")

	err := s.Write(logpipe.Entry{"message": "hi"})
	if err == nil || err.Error() != "write failed" {
		t.Fatalf("expected propagated error, got %v", err)
	}
}

func TestPrefixSink_Close(t *testing.T) {
	col := &collectingSink{}
	s := sink.NewPrefixSink(col, "message", "pre_")
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
}
