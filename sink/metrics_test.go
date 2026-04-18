package sink_test

import (
	"errors"
	"testing"

	"github.com/andybar2/logpipe"
	"github.com/andybar2/logpipe/sink"
)

func TestMetricsSink_CountsWrites(t *testing.T) {
	cs := &captureSink{}
	m := sink.NewMetricsSink(cs)

	for i := 0; i < 5; i++ {
		_ = m.Write(logpipe.Entry{Message: "ok"})
	}

	if m.Writes() != 5 {
		t.Fatalf("expected 5 writes, got %d", m.Writes())
	}
	if m.Errors() != 0 || m.Drops() != 0 {
		t.Fatal("expected no errors or drops")
	}
}

func TestMetricsSink_CountsErrors(t *testing.T) {
	es := &errorSink{err: errors.New("boom")}
	m := sink.NewMetricsSink(es)

	_ = m.Write(logpipe.Entry{Message: "fail"})

	if m.Errors() != 1 {
		t.Fatalf("expected 1 error, got %d", m.Errors())
	}
	if m.Writes() != 0 {
		t.Fatal("expected no writes")
	}
}

func TestMetricsSink_CountsDrops(t *testing.T) {
	ds := &errorSink{err: sink.ErrDropped}
	m := sink.NewMetricsSink(ds)

	_ = m.Write(logpipe.Entry{Message: "drop"})

	if m.Drops() != 1 {
		t.Fatalf("expected 1 drop, got %d", m.Drops())
	}
	if m.Writes() != 0 || m.Errors() != 0 {
		t.Fatal("unexpected writes or errors")
	}
}

func TestMetricsSink_Reset(t *testing.T) {
	cs := &captureSink{}
	m := sink.NewMetricsSink(cs)
	_ = m.Write(logpipe.Entry{Message: "x"})
	m.Reset()

	if m.Writes() != 0 || m.Errors() != 0 || m.Drops() != 0 {
		t.Fatal("expected all counters to be zero after reset")
	}
}

func TestMetricsSink_Close(t *testing.T) {
	cs := &captureSink{}
	m := sink.NewMetricsSink(cs)
	if err := m.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
}
