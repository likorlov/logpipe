package sink_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

type collectSinkAgg struct {
	mu      sync.Mutex
	entries []logpipe.Entry
	err     error
}

func (c *collectSinkAgg) Write(e logpipe.Entry) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err != nil {
		return c.err
	}
	c.entries = append(c.entries, e)
	return nil
}
func (c *collectSinkAgg) Close() error { return nil }

func TestAggregateSink_FlushOnSize(t *testing.T) {
	col := &collectSinkAgg{}
	a := sink.NewAggregateSink(col, 3, "")

	for i := 0; i < 2; i++ {
		_ = a.Write(logpipe.Entry{Level: "info", Message: "m", Fields: map[string]interface{}{"i": i}})
	}
	if len(col.entries) != 0 {
		t.Fatal("expected no flush before size reached")
	}
	_ = a.Write(logpipe.Entry{Level: "warn", Message: "last", Fields: map[string]interface{}{"x": 9}})
	if len(col.entries) != 1 {
		t.Fatalf("expected 1 merged entry, got %d", len(col.entries))
	}
	merged := col.entries[0]
	if merged.Fields["agg_count"] != 3 {
		t.Errorf("expected agg_count=3, got %v", merged.Fields["agg_count"])
	}
	if merged.Message != "last" {
		t.Errorf("expected message from last entry")
	}
}

func TestAggregateSink_FlushOnClose(t *testing.T) {
	col := &collectSinkAgg{}
	a := sink.NewAggregateSink(col, 10, "count")
	_ = a.Write(logpipe.Entry{Level: "debug", Message: "only", Fields: map[string]interface{}{"k": "v"}})
	if err := a.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry after close, got %d", len(col.entries))
	}
	if col.entries[0].Fields["count"] != 1 {
		t.Errorf("expected count=1")
	}
}

func TestAggregateSink_MergesFields(t *testing.T) {
	col := &collectSinkAgg{}
	a := sink.NewAggregateSink(col, 2, "")
	_ = a.Write(logpipe.Entry{Level: "info", Message: "a", Fields: map[string]interface{}{"a": 1, "shared": "first"}})
	_ = a.Write(logpipe.Entry{Level: "info", Message: "b", Fields: map[string]interface{}{"b": 2, "shared": "second"}})
	if col.entries[0].Fields["a"] != 1 || col.entries[0].Fields["b"] != 2 {
		t.Error("expected fields from both entries")
	}
	if col.entries[0].Fields["shared"] != "second" {
		t.Error("expected later entry to win on collision")
	}
}

func TestAggregateSink_InnerError(t *testing.T) {
	col := &collectSinkAgg{err: errors.New("fail")}
	a := sink.NewAggregateSink(col, 1, "")
	err := a.Write(logpipe.Entry{Level: "info", Message: "x", Fields: map[string]interface{}{}})
	if err == nil {
		t.Error("expected error from inner sink")
	}
}
