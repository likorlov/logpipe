package sink_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestThrottleSink_PassesFirst(t *testing.T) {
	col := &collectSink{}
	th := sink.NewThrottleSink(col, 100*time.Millisecond, nil)
	defer th.Close()

	e := logpipe.Entry{Message: "hello"}
	if err := th.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
}

func TestThrottleSink_SuppressesWithinCooldown(t *testing.T) {
	col := &collectSink{}
	th := sink.NewThrottleSink(col, 200*time.Millisecond, nil)
	defer th.Close()

	e := logpipe.Entry{Message: "burst"}
	for i := 0; i < 5; i++ {
		_ = th.Write(e)
	}
	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
}

func TestThrottleSink_AllowsAfterCooldown(t *testing.T) {
	col := &collectSink{}
	th := sink.NewThrottleSink(col, 30*time.Millisecond, nil)
	defer th.Close()

	e := logpipe.Entry{Message: "spaced"}
	_ = th.Write(e)
	time.Sleep(50 * time.Millisecond)
	_ = th.Write(e)

	if len(col.entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(col.entries))
	}
}

func TestThrottleSink_CustomKeyFn(t *testing.T) {
	col := &collectSink{}
	keyFn := func(e logpipe.Entry) string {
		if v, ok := e.Fields["id"]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
		return e.Message
	}
	th := sink.NewThrottleSink(col, 200*time.Millisecond, keyFn)
	defer th.Close()

	_ = th.Write(logpipe.Entry{Message: "msg", Fields: map[string]any{"id": "a"}})
	_ = th.Write(logpipe.Entry{Message: "msg", Fields: map[string]any{"id": "b"}})
	_ = th.Write(logpipe.Entry{Message: "msg", Fields: map[string]any{"id": "a"}})

	if len(col.entries) != 2 {
		t.Fatalf("expected 2 entries (one per unique id), got %d", len(col.entries))
	}
}
