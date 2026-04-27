package sink_test

import (
	"sync"
	"testing"
	"time"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestDebounceSink_ForwardsAfterQuiet(t *testing.T) {
	var mu sync.Mutex
	var got []logpipe.Entry
	collect := &callbackSink{fn: func(e logpipe.Entry) error {
		mu.Lock()
		got = append(got, e)
		mu.Unlock()
		return nil
	}}
	d := sink.NewDebounceSink(collect, 30*time.Millisecond)
	_ = d.Write(logpipe.Entry{"msg": "first"})
	time.Sleep(60 * time.Millisecond)
	mu.Lock()
	n := len(got)
	mu.Unlock()
	if n != 1 {
		t.Fatalf("expected 1 forwarded entry, got %d", n)
	}
	if got[0]["msg"] != "first" {
		t.Fatalf("unexpected entry: %v", got[0])
	}
	_ = d.Close()
}

func TestDebounceSink_KeepsLastOnBurst(t *testing.T) {
	var mu sync.Mutex
	var got []logpipe.Entry
	collect := &callbackSink{fn: func(e logpipe.Entry) error {
		mu.Lock()
		got = append(got, e)
		mu.Unlock()
		return nil
	}}
	d := sink.NewDebounceSink(collect, 40*time.Millisecond)
	for i := 0; i < 5; i++ {
		_ = d.Write(logpipe.Entry{"seq": i})
		time.Sleep(5 * time.Millisecond)
	}
	time.Sleep(80 * time.Millisecond)
	mu.Lock()
	n := len(got)
	last := got[len(got)-1]
	mu.Unlock()
	if n != 1 {
		t.Fatalf("expected 1 forwarded entry after burst, got %d", n)
	}
	if last["seq"] != 4 {
		t.Fatalf("expected last seq=4, got %v", last["seq"])
	}
	_ = d.Close()
}

func TestDebounceSink_FlushOnClose(t *testing.T) {
	var got []logpipe.Entry
	collect := &callbackSink{fn: func(e logpipe.Entry) error {
		got = append(got, e)
		return nil
	}}
	d := sink.NewDebounceSink(collect, 500*time.Millisecond)
	_ = d.Write(logpipe.Entry{"msg": "pending"})
	_ = d.Close()
	if len(got) != 1 {
		t.Fatalf("expected pending entry flushed on Close, got %d entries", len(got))
	}
}

func TestDebounceSink_PanicOnZeroWait(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero wait")
		}
	}()
	sink.NewDebounceSink(&callbackSink{}, 0)
}
