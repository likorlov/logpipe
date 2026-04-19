package sink_test

import (
	"errors"
	"testing"

	"github.com/yourusername/logpipe"
	"github.com/yourusername/logpipe/sink"
)

func TestTruncateSink_ShortValuePassesThrough(t *testing.T) {
	col := &collectSink{}
	tr := sink.NewTruncateSink(col, 10, "...", "msg")

	_ = tr.Write(logpipe.Entry{"msg": "hello"})
	if col.entries[0]["msg"] != "hello" {
		t.Fatalf("expected 'hello', got %v", col.entries[0]["msg"])
	}
}

func TestTruncateSink_TruncatesLongValue(t *testing.T) {
	col := &collectSink{}
	tr := sink.NewTruncateSink(col, 5, "...", "msg")

	_ = tr.Write(logpipe.Entry{"msg": "hello world"})
	got := col.entries[0]["msg"]
	if got != "hello..." {
		t.Fatalf("expected 'hello...', got %v", got)
	}
}

func TestTruncateSink_NoSuffix(t *testing.T) {
	col := &collectSink{}
	tr := sink.NewTruncateSink(col, 4, "", "msg")

	_ = tr.Write(logpipe.Entry{"msg": "truncated"})
	if col.entries[0]["msg"] != "trun" {
		t.Fatalf("expected 'trun', got %v", col.entries[0]["msg"])
	}
}

func TestTruncateSink_NonStringFieldUnchanged(t *testing.T) {
	col := &collectSink{}
	tr := sink.NewTruncateSink(col, 3, "...", "count")

	_ = tr.Write(logpipe.Entry{"count": 42})
	if col.entries[0]["count"] != 42 {
		t.Fatalf("expected 42, got %v", col.entries[0]["count"])
	}
}

func TestTruncateSink_NoMutationOfOriginal(t *testing.T) {
	col := &collectSink{}
	tr := sink.NewTruncateSink(col, 3, "", "msg")

	orig := logpipe.Entry{"msg": "hello"}
	_ = tr.Write(orig)
	if orig["msg"] != "hello" {
		t.Fatal("original entry was mutated")
	}
}

func TestTruncateSink_PropagatesError(t *testing.T) {
	errSink := &errWriteSink{err: errors.New("write failed")}
	tr := sink.NewTruncateSink(errSink, 5, "", "msg")

	err := tr.Write(logpipe.Entry{"msg": "hello world"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTruncateSink_Close(t *testing.T) {
	col := &collectSink{}
	tr := sink.NewTruncateSink(col, 10, "", "msg")
	if err := tr.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
}
