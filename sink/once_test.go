package sink_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestOnceSink_ForwardsFirst(t *testing.T) {
	col := &collectSink{}
	s := sink.NewOnceSink(col, nil)

	for i := 0; i < 5; i++ {
		if err := s.Write(logpipe.Entry{Message: "msg"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
}

func TestOnceSink_PredicateFilters(t *testing.T) {
	col := &collectSink{}
	pred := func(e logpipe.Entry) bool { return e.Level == logpipe.ERROR }
	s := sink.NewOnceSink(col, pred)

	s.Write(logpipe.Entry{Level: logpipe.INFO, Message: "info"})
	s.Write(logpipe.Entry{Level: logpipe.ERROR, Message: "err1"})
	s.Write(logpipe.Entry{Level: logpipe.ERROR, Message: "err2"})

	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
	if col.entries[0].Message != "err1" {
		t.Errorf("expected err1, got %s", col.entries[0].Message)
	}
}

func TestOnceSink_Reset(t *testing.T) {
	col := &collectSink{}
	s := sink.NewOnceSink(col, nil)

	s.Write(logpipe.Entry{Message: "first"})
	s.Reset()
	s.Write(logpipe.Entry{Message: "second"})

	if len(col.entries) != 2 {
		t.Fatalf("expected 2 entries after reset, got %d", len(col.entries))
	}
}

func TestOnceSink_ConcurrentSafe(t *testing.T) {
	col := &collectSink{}
	s := sink.NewOnceSink(col, nil)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Write(logpipe.Entry{Message: "x"})
		}()
	}
	wg.Wait()

	if len(col.entries) != 1 {
		t.Errorf("expected exactly 1 entry, got %d", len(col.entries))
	}
}

func TestOnceSink_PropagatesError(t *testing.T) {
	errSink := &errOnWriteSink{err: errors.New("boom")}
	s := sink.NewOnceSink(errSink, nil)

	if err := s.Write(logpipe.Entry{Message: "x"}); err == nil {
		t.Error("expected error from inner sink")
	}
}
