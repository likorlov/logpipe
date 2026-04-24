package sink_test

import (
	"errors"
	"testing"
	"time"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestExpireSink_ForwardsFreshEntry(t *testing.T) {
	col := &captureSink{}
	s := sink.NewExpireSink(col, 5*time.Minute)

	entry := logpipe.Entry{Fields: logpipe.Fields{
		"ts":  time.Now().Add(-1 * time.Minute),
		"msg": "hello",
	}}

	if err := s.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
}

func TestExpireSink_DropsExpiredEntry(t *testing.T) {
	col := &captureSink{}
	s := sink.NewExpireSink(col, 5*time.Minute)

	entry := logpipe.Entry{Fields: logpipe.Fields{
		"ts":  time.Now().Add(-10 * time.Minute),
		"msg": "stale",
	}}

	if err := s.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(col.entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(col.entries))
	}
}

func TestExpireSink_MissingFieldPassesThrough(t *testing.T) {
	col := &captureSink{}
	s := sink.NewExpireSink(col, 5*time.Minute)

	entry := logpipe.Entry{Fields: logpipe.Fields{"msg": "no timestamp"}}

	if err := s.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
}

func TestExpireSink_NonTimeFieldPassesThrough(t *testing.T) {
	col := &captureSink{}
	s := sink.NewExpireSink(col, 5*time.Minute)

	entry := logpipe.Entry{Fields: logpipe.Fields{"ts": "not-a-time"}}

	if err := s.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
}

func TestExpireSink_CustomField(t *testing.T) {
	col := &captureSink{}
	s := sink.NewExpireSink(col, time.Minute, sink.WithExpireField("created_at"))

	entry := logpipe.Entry{Fields: logpipe.Fields{
		"created_at": time.Now().Add(-2 * time.Minute),
	}}

	if err := s.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(col.entries) != 0 {
		t.Fatalf("expected 0 entries (expired), got %d", len(col.entries))
	}
}

func TestExpireSink_PropagatesInnerError(t *testing.T) {
	want := errors.New("inner error")
	s := sink.NewExpireSink(&errSink{err: want}, time.Hour)

	entry := logpipe.Entry{Fields: logpipe.Fields{"ts": time.Now()}}

	if got := s.Write(entry); !errors.Is(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestExpireSink_Close(t *testing.T) {
	col := &captureSink{}
	s := sink.NewExpireSink(col, time.Minute)
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
}
