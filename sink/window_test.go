package sink_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestWindowSink_ForwardsEntry(t *testing.T) {
	col := &collectSink{}
	w := sink.NewWindowSink(col, 5*time.Second)

	entry := logpipe.Entry{"msg": "hello"}
	if err := w.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(col.entries) != 1 {
		t.Fatalf("expected 1 forwarded entry, got %d", len(col.entries))
	}
}

func TestWindowSink_EntriesWithinWindow(t *testing.T) {
	col := &collectSink{}
	w := sink.NewWindowSink(col, 5*time.Second)

	w.Write(logpipe.Entry{"msg": "a"})
	w.Write(logpipe.Entry{"msg": "b"})

	entries := w.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries in window, got %d", len(entries))
	}
}

func TestWindowSink_EvictsExpired(t *testing.T) {
	col := &collectSink{}
	w := sink.NewWindowSink(col, 50*time.Millisecond)

	w.Write(logpipe.Entry{"msg": "old"})
	time.Sleep(80 * time.Millisecond)
	w.Write(logpipe.Entry{"msg": "new"})

	entries := w.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after eviction, got %d", len(entries))
	}
	if entries[0]["msg"] != "new" {
		t.Fatalf("expected 'new', got %v", entries[0]["msg"])
	}
}

func TestWindowSink_Close(t *testing.T) {
	col := &collectSink{}
	w := sink.NewWindowSink(col, time.Second)
	if err := w.Close(); err != nil {
		t.Fatalf("unexpected error on close: %v", err)
	}
	if !col.closed {
		t.Fatal("expected inner sink to be closed")
	}
}
