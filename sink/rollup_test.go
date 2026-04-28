package sink_test

import (
	"sync"
	"testing"
	"time"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

// captureSink collects every entry written to it.
type captureSink struct {
	mu      sync.Mutex
	entries []logpipe.Entry
}

func (c *captureSink) Write(e logpipe.Entry) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = append(c.entries, e)
	return nil
}

func (c *captureSink) Close() error { return nil }

func (c *captureSink) all() []logpipe.Entry {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]logpipe.Entry, len(c.entries))
	copy(out, c.entries)
	return out
}

func TestRollupSink_FlushOnClose(t *testing.T) {
	cap := &captureSink{}
	s := sink.NewRollupSink(cap, "latency", 10*time.Second)

	for _, v := range []float64{10, 20, 30} {
		_ = s.Write(logpipe.Entry{Fields: map[string]interface{}{"latency": v}})
	}
	_ = s.Close()

	entries := cap.all()
	if len(entries) != 1 {
		t.Fatalf("expected 1 rollup entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Fields["count"].(int64) != 3 {
		t.Errorf("count: want 3, got %v", e.Fields["count"])
	}
	if e.Fields["sum"].(float64) != 60 {
		t.Errorf("sum: want 60, got %v", e.Fields["sum"])
	}
	if e.Fields["min"].(float64) != 10 {
		t.Errorf("min: want 10, got %v", e.Fields["min"])
	}
	if e.Fields["max"].(float64) != 30 {
		t.Errorf("max: want 30, got %v", e.Fields["max"])
	}
	if e.Fields["field"].(string) != "latency" {
		t.Errorf("field: want latency, got %v", e.Fields["field"])
	}
}

func TestRollupSink_SkipsNonNumericField(t *testing.T) {
	cap := &captureSink{}
	s := sink.NewRollupSink(cap, "latency", 10*time.Second)
	_ = s.Write(logpipe.Entry{Fields: map[string]interface{}{"latency": "fast"}})
	_ = s.Close()
	if len(cap.all()) != 0 {
		t.Error("expected no rollup entry for non-numeric field")
	}
}

func TestRollupSink_SkipsMissingField(t *testing.T) {
	cap := &captureSink{}
	s := sink.NewRollupSink(cap, "latency", 10*time.Second)
	_ = s.Write(logpipe.Entry{Fields: map[string]interface{}{"msg": "hello"}})
	_ = s.Close()
	if len(cap.all()) != 0 {
		t.Error("expected no rollup entry when field absent")
	}
}

func TestRollupSink_FlushOnInterval(t *testing.T) {
	cap := &captureSink{}
	s := sink.NewRollupSink(cap, "bytes", 50*time.Millisecond)
	_ = s.Write(logpipe.Entry{Fields: map[string]interface{}{"bytes": 100.0}})
	time.Sleep(120 * time.Millisecond)
	_ = s.Close()
	if len(cap.all()) == 0 {
		t.Error("expected at least one rollup flush from ticker")
	}
}

func TestRollupSink_PanicOnZeroWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for zero window")
		}
	}()
	sink.NewRollupSink(&captureSink{}, "x", 0)
}
