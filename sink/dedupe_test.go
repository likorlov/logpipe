package sink_test

import (
	"sync"
	"testing"
	"time"

	"github.com/your/logpipe"
	"github.com/your/logpipe/sink"
)

type countSink struct {
	mu    sync.Mutex
	count int
}

func (c *countSink) Write(_ logpipe.Entry) error {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
	return nil
}
func (c *countSink) Close() error { return nil }
func (c *countSink) n() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func TestDedupeSink_PassesUnique(t *testing.T) {
	cs := &countSink{}
	d := sink.NewDedupeSink(cs, time.Second)

	_ = d.Write(logpipe.Entry{Level: logpipe.Info, Message: "hello"})
	_ = d.Write(logpipe.Entry{Level: logpipe.Info, Message: "world"})

	if cs.n() != 2 {
		t.Fatalf("expected 2 writes, got %d", cs.n())
	}
}

func TestDedupeSink_SuppressDuplicate(t *testing.T) {
	cs := &countSink{}
	d := sink.NewDedupeSink(cs, time.Second)

	for i := 0; i < 5; i++ {
		_ = d.Write(logpipe.Entry{Level: logpipe.Error, Message: "boom"})
	}

	if cs.n() != 1 {
		t.Fatalf("expected 1 write, got %d", cs.n())
	}
}

func TestDedupeSink_AllowsAfterWindow(t *testing.T) {
	cs := &countSink{}
	d := sink.NewDedupeSink(cs, 20*time.Millisecond)

	_ = d.Write(logpipe.Entry{Level: logpipe.Warn, Message: "transient"})
	time.Sleep(30 * time.Millisecond)
	_ = d.Write(logpipe.Entry{Level: logpipe.Warn, Message: "transient"})

	if cs.n() != 2 {
		t.Fatalf("expected 2 writes after window expiry, got %d", cs.n())
	}
}

func TestDedupeSink_DifferentLevelsSameMessage(t *testing.T) {
	cs := &countSink{}
	d := sink.NewDedupeSink(cs, time.Second)

	_ = d.Write(logpipe.Entry{Level: logpipe.Info, Message: "msg"})
	_ = d.Write(logpipe.Entry{Level: logpipe.Error, Message: "msg"})

	if cs.n() != 2 {
		t.Fatalf("expected 2 writes for different levels, got %d", cs.n())
	}
}
