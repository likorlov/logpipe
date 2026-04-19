package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestReplaySink_RecordsEntries(t *testing.T) {
	r := sink.NewReplaySink(nil)
	e1 := logpipe.Entry{Message: "one"}
	e2 := logpipe.Entry{Message: "two"}
	_ = r.Write(e1)
	_ = r.Write(e2)

	entries := r.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Message != "one" || entries[1].Message != "two" {
		t.Fatalf("unexpected entries: %+v", entries)
	}
}

func TestReplaySink_Replay(t *testing.T) {
	r := sink.NewReplaySink(nil)
	_ = r.Write(logpipe.Entry{Message: "a"})
	_ = r.Write(logpipe.Entry{Message: "b"})

	dst := sink.NewSnapshotSink(10)
	if err := r.Replay(dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := dst.Entries(); len(got) != 2 {
		t.Fatalf("expected 2 entries replayed, got %d", len(got))
	}
}

func TestReplaySink_Reset(t *testing.T) {
	r := sink.NewReplaySink(nil)
	_ = r.Write(logpipe.Entry{Message: "x"})
	r.Reset()
	if len(r.Entries()) != 0 {
		t.Fatal("expected empty after reset")
	}
}

func TestReplaySink_ForwardsToInner(t *testing.T) {
	inner := sink.NewSnapshotSink(10)
	r := sink.NewReplaySink(inner)
	_ = r.Write(logpipe.Entry{Message: "fwd"})

	if len(inner.Entries()) != 1 {
		t.Fatal("expected entry forwarded to inner sink")
	}
}

func TestReplaySink_ReplayPropagatesError(t *testing.T) {
	r := sink.NewReplaySink(nil)
	_ = r.Write(logpipe.Entry{Message: "err"})

	boom := errors.New("boom")
	errSink := &errWriteSink{err: boom}
	if err := r.Replay(errSink); !errors.Is(err, boom) {
		t.Fatalf("expected boom, got %v", err)
	}
}

type errWriteSink struct{ err error }

func (e *errWriteSink) Write(_ logpipe.Entry) error { return e.err }
func (e *errWriteSink) Close() error                { return nil }
