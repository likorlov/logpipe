package sink_test

import (
	"strconv"
	"testing"

	"github.com/your/logpipe"
	"github.com/your/logpipe/sink"
)

func TestSequenceSink_InjectsField(t *testing.T) {
	col := &collectSink{}
	s := sink.NewSequenceSink(col, "seq")

	for i := 0; i < 3; i++ {
		if err := s.Write(logpipe.Entry{Message: "hello", Fields: map[string]any{}}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if len(col.entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(col.entries))
	}
	for i, e := range col.entries {
		want := strconv.Itoa(i + 1)
		got, ok := e.Fields["seq"].(string)
		if !ok || got != want {
			t.Errorf("entry %d: seq=%q want %q", i, got, want)
		}
	}
}

func TestSequenceSink_DefaultField(t *testing.T) {
	col := &collectSink{}
	s := sink.NewSequenceSink(col, "")
	_ = s.Write(logpipe.Entry{Message: "x", Fields: map[string]any{}})
	if _, ok := col.entries[0].Fields["seq"]; !ok {
		t.Error("expected default field 'seq'")
	}
}

func TestSequenceSink_NoMutationOfOriginal(t *testing.T) {
	col := &collectSink{}
	s := sink.NewSequenceSink(col, "seq")
	orig := logpipe.Entry{Message: "m", Fields: map[string]any{"a": 1}}
	_ = s.Write(orig)
	if _, ok := orig.Fields["seq"]; ok {
		t.Error("original entry was mutated")
	}
}

func TestSequenceSink_Counter(t *testing.T) {
	col := &collectSink{}
	s := sink.NewSequenceSink(col, "seq")
	for i := 0; i < 5; i++ {
		_ = s.Write(logpipe.Entry{Message: "x", Fields: map[string]any{}})
	}
	if s.Counter() != 5 {
		t.Errorf("counter=%d want 5", s.Counter())
	}
}

func TestSequenceSink_Close(t *testing.T) {
	col := &collectSink{}
	s := sink.NewSequenceSink(col, "seq")
	if err := s.Close(); err != nil {
		t.Errorf("unexpected close error: %v", err)
	}
}
