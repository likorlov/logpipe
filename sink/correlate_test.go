package sink_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/your-org/logpipe"
	"github.com/your-org/logpipe/sink"
)

func TestCorrelateSink_InjectsField(t *testing.T) {
	col := &collectSink{}
	s := sink.NewCorrelateSink(col, "", func() string { return "abc-123" })

	if err := s.Write(logpipe.Entry{"msg": "hello"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := col.entries[0]
	if got["correlation_id"] != "abc-123" {
		t.Errorf("expected correlation_id=abc-123, got %v", got["correlation_id"])
	}
}

func TestCorrelateSink_CustomField(t *testing.T) {
	col := &collectSink{}
	s := sink.NewCorrelateSink(col, "trace_id", func() string { return "t-1" })

	_ = s.Write(logpipe.Entry{})
	if col.entries[0]["trace_id"] != "t-1" {
		t.Errorf("expected trace_id=t-1, got %v", col.entries[0]["trace_id"])
	}
}

func TestCorrelateSink_EntryFieldWins(t *testing.T) {
	col := &collectSink{}
	s := sink.NewCorrelateSink(col, "correlation_id", func() string { return "auto" })

	_ = s.Write(logpipe.Entry{"correlation_id": "manual"})
	if col.entries[0]["correlation_id"] != "manual" {
		t.Errorf("expected entry field to win, got %v", col.entries[0]["correlation_id"])
	}
}

func TestCorrelateSink_Rotate(t *testing.T) {
	col := &collectSink{}
	n := 0
	s := sink.NewCorrelateSink(col, "", func() string {
		n++
		return fmt.Sprintf("id-%d", n)
	})

	_ = s.Write(logpipe.Entry{})
	first := col.entries[0]["correlation_id"]

	s.Rotate()
	_ = s.Write(logpipe.Entry{})
	second := col.entries[1]["correlation_id"]

	if first == second {
		t.Errorf("expected different IDs after Rotate, got %v both times", first)
	}
}

func TestCorrelateSink_NoMutationOfOriginal(t *testing.T) {
	col := &collectSink{}
	s := sink.NewCorrelateSink(col, "correlation_id", func() string { return "x" })

	orig := logpipe.Entry{"msg": "hi"}
	_ = s.Write(orig)

	if _, ok := orig["correlation_id"]; ok {
		t.Error("original entry was mutated")
	}
}

func TestCorrelateSink_ConcurrentSafe(t *testing.T) {
	col := &collectSink{}
	s := sink.NewCorrelateSink(col, "", func() string { return "c" })

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.Write(logpipe.Entry{"x": 1})
			s.Rotate()
			_ = s.Current()
		}()
	}
	wg.Wait()
}

func TestCorrelateSink_Close(t *testing.T) {
	col := &collectSink{}
	s := sink.NewCorrelateSink(col, "", nil)
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
}
