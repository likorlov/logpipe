package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

// pipelineCapture records the last entry it received and optionally returns an error.
type pipelineCapture struct {
	last   logpipe.Entry
	failOn int // return error on the n-th Write (1-based); 0 = never
	calls  int
	closed bool
}

func (c *pipelineCapture) Write(e logpipe.Entry) error {
	c.calls++
	c.last = e
	if c.failOn > 0 && c.calls >= c.failOn {
		return errors.New("stage error")
	}
	return nil
}

func (c *pipelineCapture) Close() error {
	c.closed = true
	return nil
}

func TestPipelineSink_ForwardsEntry(t *testing.T) {
	a := &pipelineCapture{}
	b := &pipelineCapture{}
	p := sink.NewPipelineSink(a, b)

	entry := logpipe.Entry{"msg": "hello"}
	if err := p.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.calls != 1 || b.calls != 1 {
		t.Fatalf("expected both stages called once, got a=%d b=%d", a.calls, b.calls)
	}
}

func TestPipelineSink_HaltsOnError(t *testing.T) {
	a := &pipelineCapture{failOn: 1}
	b := &pipelineCapture{}
	p := sink.NewPipelineSink(a, b)

	if err := p.Write(logpipe.Entry{"msg": "x"}); err == nil {
		t.Fatal("expected error from stage a")
	}
	if b.calls != 0 {
		t.Fatalf("expected stage b to be skipped, got %d calls", b.calls)
	}
}

func TestPipelineSink_Close(t *testing.T) {
	a := &pipelineCapture{}
	b := &pipelineCapture{}
	p := sink.NewPipelineSink(a, b)

	if err := p.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
	if !a.closed || !b.closed {
		t.Fatal("expected all stages to be closed")
	}
}

func TestPipelineSink_PanicOnEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for empty stages")
		}
	}()
	sink.NewPipelineSink()
}
