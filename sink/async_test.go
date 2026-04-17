package sink_test

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

type recordSink struct {
	entries []logpipe.Entry
	failWith error
}

func (r *recordSink) Write(e logpipe.Entry) error {
	if r.failWith != nil {
		return r.failWith
	}
	r.entries = append(r.entries,
}
func (r *recordSink) Close() error { return nil }

func TestAsyncSink_DeliversEntries(t *testing.T) {
	rec := &recordSink{}
	as := sink.NewAsyncSink(rec, 16)

	entry := logpipe.Entry{Level: logpipe.INFO, Message: "hello async"}
	if err := as.Write(entry); err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}
	if err := as.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
	if len(rec.entries) != 1 || rec.entries[0].Message != "hello async" {
		t.Fatalf("expected 1 entry, got %v", rec.entries)
	}
}

func TestAsyncSink_QueueFull(t *testing.T) {
	// Use a sink that blocks so the queue fills up.
	block := make(chan struct{})
	blockedSink := &blockingSink{block: block}
	as := sink.NewAsyncSink(blockedSink, 1)

	_ = as.Write(logpipe.Entry{Message: "first"})
	_ = as.Write(logpipe.Entry{Message: "second"}) // fills queue
	err := as.Write(logpipe.Entry{Message: "overflow"})
	if err == nil {
		t.Fatal("expected queue-full error")
	}
	close(block)
	_ = as.Close()
}

func TestAsyncSink_ErrFuncCalled(t *testing.T) {
	var count int64
	rec := &recordSink{failWith: errors.New("sink down")}
	as := sink.NewAsyncSink(rec, 8)
	as.ErrFunc = func(err error) { atomic.AddInt64(&count, 1) }

	_ = as.Write(logpipe.Entry{Message: "boom"})
	_ = as.Close()

	time.Sleep(10 * time.Millisecond)
	if atomic.LoadInt64(&count) == 0 {
		t.Fatal("expected ErrFunc to be called")
	}
}

type blockingSink struct{ block chan struct{} }

func (b *blockingSink) Write(e logpipe.Entry) error { <-b.block; return nil }
func (b *blockingSink) Close() error                { return nil }
