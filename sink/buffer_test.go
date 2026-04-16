package sink_test

import (
	"sync"
	"testing"
	"time"

	"github.com/example/logpipe"
	"github.com/example/logpipe/sink"
)

type captureSink struct {
	mu      sync.Mutex
	entries []logpipe.Entry
	closed  bool
}

func (c *captureSink) Write(e logpipe.Entry) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = append(c.entries, e)
	return nil
}

func (c *captureSink) Close() error {
	c.closed = true
	return nil
}

func (c *captureSink) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}

func TestBufferedSink_FlushOnSize(t *testing.T) {
	inner := &captureSink{}
	b := sink.NewBufferedSink(inner, 3, 10*time.Second)
	defer b.Close()

	for i := 0; i < 3; i++ {
		_ = b.Write(logpipe.Entry{Message: "msg"})
	}
	time.Sleep(20 * time.Millisecond)
	if inner.Len() != 3 {
		t.Fatalf("expected 3 entries flushed, got %d", inner.Len())
	}
}

func TestBufferedSink_FlushOnInterval(t *testing.T) {
	inner := &captureSink{}
	b := sink.NewBufferedSink(inner, 100, 50*time.Millisecond)
	defer b.Close()

	_ = b.Write(logpipe.Entry{Message: "hello"})
	time.Sleep(120 * time.Millisecond)
	if inner.Len() != 1 {
		t.Fatalf("expected 1 entry after interval flush, got %d", inner.Len())
	}
}

func TestBufferedSink_FlushOnClose(t *testing.T) {
	inner := &captureSink{}
	b := sink.NewBufferedSink(inner, 100, 10*time.Second)

	_ = b.Write(logpipe.Entry{Message: "close-flush"})
	_ = b.Close()

	if inner.Len() != 1 {
		t.Fatalf("expected 1 entry flushed on close, got %d", inner.Len())
	}
	if !inner.closed {
		t.Fatal("expected inner sink to be closed")
	}
}
