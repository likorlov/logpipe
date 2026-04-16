package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

// captureSink records every entry written to it.
type captureSink struct {
	entries []logpipe.Entry
	closed  bool
}

func (c *captureSink) Write(e logpipe.Entry) error { c.entries = append(c.entries, e); return nil }
func (c *captureSink) Close() error                { c.closed = true; return nil }

// errorSink always returns an error on Write.
type errorSink struct{}

func (e *errorSink) Write(_ logpipe.Entry) error { return errors.New("write error") }
func (e *errorSink) Close() error                { return nil }

func TestFilterSink_PassesMatchingEntries(t *testing.T) {
	cap := &captureSink{}
	fs := sink.NewFilterSink(cap, sink.LevelFilter(logpipe.LevelWarn))

	_ = fs.Write(logpipe.Entry{Level: logpipe.LevelDebug, Message: "debug"})
	_ = fs.Write(logpipe.Entry{Level: logpipe.LevelInfo, Message: "info"})
	_ = fs.Write(logpipe.Entry{Level: logpipe.LevelWarn, Message: "warn"})
	_ = fs.Write(logpipe.Entry{Level: logpipe.LevelError, Message: "error"})

	if len(cap.entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(cap.entries))
	}
	if cap.entries[0].Message != "warn" {
		t.Errorf("unexpected first entry: %s", cap.entries[0].Message)
	}
}

func TestFilterSink_DropsNonMatchingEntries(t *testing.T) {
	cap := &captureSink{}
	fs := sink.NewFilterSink(cap, sink.LevelFilter(logpipe.LevelError))

	_ = fs.Write(logpipe.Entry{Level: logpipe.LevelInfo, Message: "info"})

	if len(cap.entries) != 0 {
		t.Errorf("expected no entries, got %d", len(cap.entries))
	}
}

func TestFilterSink_FieldFilter(t *testing.T) {
	cap := &captureSink{}
	fs := sink.NewFilterSink(cap, sink.FieldFilter("request_id"))

	_ = fs.Write(logpipe.Entry{Message: "no field"})
	_ = fs.Write(logpipe.Entry{Message: "has field", Fields: map[string]any{"request_id": "abc"}})

	if len(cap.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(cap.entries))
	}
	if cap.entries[0].Message != "has field" {
		t.Errorf("unexpected entry: %s", cap.entries[0].Message)
	}
}

func TestFilterSink_Close(t *testing.T) {
	cap := &captureSink{}
	fs := sink.NewFilterSink(cap, sink.LevelFilter(logpipe.LevelInfo))
	_ = fs.Close()
	if !cap.closed {
		t.Error("expected inner sink to be closed")
	}
}

func TestFilterSink_PropagatesWriteError(t *testing.T) {
	fs := sink.NewFilterSink(&errorSink{}, sink.LevelFilter(logpipe.LevelDebug))
	err := fs.Write(logpipe.Entry{Level: logpipe.LevelInfo, Message: "x"})
	if err == nil {
		t.Error("expected error from inner sink")
	}
}
