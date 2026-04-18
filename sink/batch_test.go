package sink_test

import (
	"sync"
	"testing"
	"time"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func collectBatch(t *testing.T) (func([]logpipe.Entry) error, func() []logpipe.Entry) {
	t.Helper()
	var mu sync.Mutex
	var all []logpipe.Entry
	fn := func(batch []logpipe.Entry) error {
		mu.Lock()
		all = append(all, batch...)
		mu.Unlock()
		return nil
	}
	get := func() []logpipe.Entry {
		mu.Lock()
		defer mu.Unlock()
		out := make([]logpipe.Entry, len(all))
		copy(out, all)
		return out
	}
	return fn, get
}

func TestBatchSink_FlushOnSize(t *testing.T) {
	flushFn, get := collectBatch(t)
	s := sink.NewBatchSink(3, 10*time.Second, flushFn)
	defer s.Close()

	for i := 0; i < 3; i++ {
		if err := s.Write(logpipe.Entry{Message: "msg"}); err != nil {
			t.Fatalf("write %d: %v", i, err)
		}
	}
	time.Sleep(20 * time.Millisecond)
	if got := len(get()); got != 3 {
		t.Fatalf("expected 3 entries flushed, got %d", got)
	}
}

func TestBatchSink_FlushOnInterval(t *testing.T) {
	flushFn, get := collectBatch(t)
	s := sink.NewBatchSink(100, 50*time.Millisecond, flushFn)
	defer s.Close()

	_ = s.Write(logpipe.Entry{Message: "tick"})
	time.Sleep(120 * time.Millisecond)
	if got := len(get()); got != 1 {
		t.Fatalf("expected 1 entry after interval flush, got %d", got)
	}
}

func TestBatchSink_FlushOnClose(t *testing.T) {
	flushFn, get := collectBatch(t)
	s := sink.NewBatchSink(100, 10*time.Second, flushFn)

	_ = s.Write(logpipe.Entry{Message: "close"})
	if err := s.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	if got := len(get()); got != 1 {
		t.Fatalf("expected 1 entry flushed on close, got %d", got)
	}
}
