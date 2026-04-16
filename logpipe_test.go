package logpipe_test

import (
	"errors"
	"testing"

	"github.com/example/logpipe"
)

// recordSink captures written entries for assertions.
type recordSink struct {
	entries []logpipe.Entry
	fail    bool
}

func (r *recordSink) Write(e logpipe.Entry) error {
	if r.fail {
		return errors.New("sink error")
	}
	r.entries = append(r.entries, e)
	return nil
}
func (r *recordSink) Close() error { return nil }

func TestLogger_FiltersLevel(t *testing.T) {
	l := logpipe.New(logpipe.WARN)
	rec := &recordSink{}
	l.AddSink(rec)

	_ = l.Log(logpipe.DEBUG, "ignored", nil)
	_ = l.Log(logpipe.INFO, "also ignored", nil)
	_ = l.Log(logpipe.WARN, "kept", nil)

	if len(rec.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(rec.entries))
	}
	if rec.entries[0].Message != "kept" {
		t.Errorf("unexpected message: %s", rec.entries[0].Message)
	}
}

func TestLogger_FanOut(t *testing.T) {
	l := logpipe.New(logpipe.DEBUG)
	a, b := &recordSink{}, &recordSink{}
	l.AddSink(a)
	l.AddSink(b)

	_ = l.Log(logpipe.INFO, "broadcast", map[string]any{"x": 1})

	if len(a.entries) != 1 || len(b.entries) != 1 {
		t.Error("expected both sinks to receive the entry")
	}
}

func TestLogger_SinkError(t *testing.T) {
	l := logpipe.New(logpipe.DEBUG)
	l.AddSink(&recordSink{fail: true})
	if err := l.Log(logpipe.INFO, "msg", nil); err == nil {
		t.Error("expected error from failing sink")
	}
}
