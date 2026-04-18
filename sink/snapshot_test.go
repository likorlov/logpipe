package sink_test

import (
	"testing"
	"time"

	"github.com/your/logpipe"
	"github.com/your/logpipe/sink"
)

func TestSnapshotSink_StoresEntries(t *testing.T) {
	s := sink.NewSnapshotSink(10)
	e := logpipe.Entry{Level: logpipe.INFO, Message: "hello", Time: time.Now()}
	if err := s.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", s.Len())
	}
	if s.Entries()[0].Message != "hello" {
		t.Fatalf("unexpected message: %s", s.Entries()[0].Message)
	}
}

func TestSnapshotSink_EvictsOldest(t *testing.T) {
	s := sink.NewSnapshotSink(3)
	for i, msg := range []string{"a", "b", "c", "d"} {
		_ = i
		s.Write(logpipe.Entry{Message: msg, Time: time.Now()})
	}
	entries := s.Entries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Message != "b" {
		t.Fatalf("expected oldest to be evicted, got %s", entries[0].Message)
	}
}

func TestSnapshotSink_Reset(t *testing.T) {
	s := sink.NewSnapshotSink(10)
	s.Write(logpipe.Entry{Message: "x", Time: time.Now()})
	s.Reset()
	if s.Len() != 0 {
		t.Fatalf("expected 0 after reset, got %d", s.Len())
	}
}

func TestSnapshotSink_EntriesIsCopy(t *testing.T) {
	s := sink.NewSnapshotSink(10)
	s.Write(logpipe.Entry{Message: "orig", Time: time.Now()})
	copy1 := s.Entries()
	copy1[0].Message = "mutated"
	if s.Entries()[0].Message != "orig" {
		t.Fatal("Entries should return a copy")
	}
}

func TestSnapshotSink_Close(t *testing.T) {
	s := sink.NewSnapshotSink(5)
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected error on close: %v", err)
	}
}
