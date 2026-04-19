package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestCounterSink_CountsWrites(t *testing.T) {
	inner := &captureSink{}
	c := sink.NewCounterSink(inner)

	for i := 0; i < 5; i++ {
		_ = c.Write(logpipe.Entry{Message: "ok"})
	}

	if c.Writes() != 5 {
		t.Fatalf("expected 5 writes, got %d", c.Writes())
	}
	if c.Drops() != 0 {
		t.Fatalf("expected 0 drops, got %d", c.Drops())
	}
	if c.Total() != 5 {
		t.Fatalf("expected total 5, got %d", c.Total())
	}
}

func TestCounterSink_CountsDrops(t *testing.T) {
	errSink := &errorSink{err: errors.New("boom")}
	c := sink.NewCounterSink(errSink)

	for i := 0; i < 3; i++ {
		_ = c.Write(logpipe.Entry{Message: "fail"})
	}

	if c.Drops() != 3 {
		t.Fatalf("expected 3 drops, got %d", c.Drops())
	}
	if c.Writes() != 0 {
		t.Fatalf("expected 0 writes, got %d", c.Writes())
	}
}

func TestCounterSink_Reset(t *testing.T) {
	inner := &captureSink{}
	c := sink.NewCounterSink(inner)
	_ = c.Write(logpipe.Entry{Message: "a"})
	_ = c.Write(logpipe.Entry{Message: "b"})

	c.Reset()

	if c.Writes() != 0 || c.Drops() != 0 {
		t.Fatal("expected counters to be zero after Reset")
	}
}

func TestCounterSink_Close(t *testing.T) {
	inner := &captureSink{}
	c := sink.NewCounterSink(inner)
	if err := c.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
}
