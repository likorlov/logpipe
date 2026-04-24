package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestWatermarkSink_HighPassesAboveThreshold(t *testing.T) {
	col := &collectSink{}
	w := sink.NewWatermarkSink(col, "latency", 100.0)

	_ = w.Write(logpipe.Entry{Fields: map[string]interface{}{"latency": 150.0}})

	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
}

func TestWatermarkSink_HighDropsBelowThreshold(t *testing.T) {
	col := &collectSink{}
	w := sink.NewWatermarkSink(col, "latency", 100.0)

	_ = w.Write(logpipe.Entry{Fields: map[string]interface{}{"latency": 50.0}})

	if len(col.entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(col.entries))
	}
}

func TestWatermarkSink_HighPassesAtThreshold(t *testing.T) {
	col := &collectSink{}
	w := sink.NewWatermarkSink(col, "latency", 100.0)

	_ = w.Write(logpipe.Entry{Fields: map[string]interface{}{"latency": 100.0}})

	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
}

func TestWatermarkSink_LowPassesBelowThreshold(t *testing.T) {
	col := &collectSink{}
	w := sink.NewWatermarkSink(col, "score", 50.0, sink.WatermarkLow())

	_ = w.Write(logpipe.Entry{Fields: map[string]interface{}{"score": 30.0}})

	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
}

func TestWatermarkSink_LowDropsAboveThreshold(t *testing.T) {
	col := &collectSink{}
	w := sink.NewWatermarkSink(col, "score", 50.0, sink.WatermarkLow())

	_ = w.Write(logpipe.Entry{Fields: map[string]interface{}{"score": 80.0}})

	if len(col.entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(col.entries))
	}
}

func TestWatermarkSink_MissingFieldPassesThrough(t *testing.T) {
	col := &collectSink{}
	w := sink.NewWatermarkSink(col, "latency", 100.0)

	_ = w.Write(logpipe.Entry{Fields: map[string]interface{}{"msg": "hello"}})

	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
}

func TestWatermarkSink_NonNumericFieldPassesThrough(t *testing.T) {
	col := &collectSink{}
	w := sink.NewWatermarkSink(col, "latency", 100.0)

	_ = w.Write(logpipe.Entry{Fields: map[string]interface{}{"latency": "fast"}})

	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
}

func TestWatermarkSink_PropagatesError(t *testing.T) {
	sentinel := errors.New("write failed")
	fs := &failSink{err: sentinel}
	w := sink.NewWatermarkSink(fs, "v", 0.0)

	err := w.Write(logpipe.Entry{Fields: map[string]interface{}{"v": 1.0}})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestWatermarkSink_Close(t *testing.T) {
	col := &collectSink{}
	w := sink.NewWatermarkSink(col, "v", 0.0)
	if err := w.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
}
