package sink_test

import (
	"errors"
	"testing"

	"github.com/your-org/logpipe"
	"github.com/your-org/logpipe/sink"
)

func TestJournalSink_InjectsIndex(t *testing.T) {
	col := &collectSink{}
	j := sink.NewJournalSink(col, "")

	if err := j.Write(logpipe.Entry{"msg": "hello"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
	got, ok := col.entries[0]["_journal"]
	if !ok {
		t.Fatal("expected _journal field to be set")
	}
	if got.(uint64) != 1 {
		t.Errorf("expected index 1, got %v", got)
	}
}

func TestJournalSink_CustomField(t *testing.T) {
	col := &collectSink{}
	j := sink.NewJournalSink(col, "seq")

	_ = j.Write(logpipe.Entry{"msg": "a"})
	_ = j.Write(logpipe.Entry{"msg": "b"})

	if col.entries[1]["seq"].(uint64) != 2 {
		t.Errorf("expected seq=2 on second entry, got %v", col.entries[1]["seq"])
	}
}

func TestJournalSink_MonotonicallyIncreasing(t *testing.T) {
	col := &collectSink{}
	j := sink.NewJournalSink(col, "")

	for i := 0; i < 5; i++ {
		_ = j.Write(logpipe.Entry{"i": i})
	}

	for i, e := range col.entries {
		want := uint64(i + 1)
		got := e["_journal"].(uint64)
		if got != want {
			t.Errorf("entry %d: expected _journal=%d, got %d", i, want, got)
		}
	}
}

func TestJournalSink_Reset(t *testing.T) {
	col := &collectSink{}
	j := sink.NewJournalSink(col, "")

	_ = j.Write(logpipe.Entry{"msg": "x"})
	_ = j.Write(logpipe.Entry{"msg": "y"})
	if j.Index() != 2 {
		t.Fatalf("expected index 2, got %d", j.Index())
	}

	j.Reset()
	if j.Index() != 0 {
		t.Errorf("expected index 0 after reset, got %d", j.Index())
	}

	_ = j.Write(logpipe.Entry{"msg": "z"})
	if j.Index() != 1 {
		t.Errorf("expected index 1 after reset+write, got %d", j.Index())
	}
}

func TestJournalSink_NoMutationOfOriginal(t *testing.T) {
	col := &collectSink{}
	j := sink.NewJournalSink(col, "")

	orig := logpipe.Entry{"msg": "hello"}
	_ = j.Write(orig)

	if _, ok := orig["_journal"]; ok {
		t.Error("original entry was mutated")
	}
}

func TestJournalSink_PropagatesError(t *testing.T) {
	errSink := &errOnWriteSink{err: errors.New("boom")}
	j := sink.NewJournalSink(errSink, "")

	err := j.Write(logpipe.Entry{"msg": "hi"})
	if err == nil || err.Error() != "boom" {
		t.Errorf("expected boom error, got %v", err)
	}
}

func TestJournalSink_Close(t *testing.T) {
	col := &collectSink{}
	j := sink.NewJournalSink(col, "")
	if err := j.Close(); err != nil {
		t.Errorf("unexpected close error: %v", err)
	}
}
