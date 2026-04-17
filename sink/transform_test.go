package sink_test

import (
	"errors"
	"testing"
	"time"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestTransformSink_ModifiesEntry(t *testing.T) {
	var got logpipe.Entry
	cap := &captureSink{}
	ts := sink.NewTransformSink(cap, func(e logpipe.Entry) (logpipe.Entry, bool) {
		e.Fields["added"] = "yes"
		return e, true
	})

	e := logpipe.Entry{Level: logpipe.INFO, Message: "hello", Time: time.Now(), Fields: map[string]interface{}{}}
	if err := ts.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got = cap.entries[0]
	if got.Fields["added"] != "yes" {
		t.Errorf("expected added field, got %v", got.Fields)
	}
}

func TestTransformSink_DropsEntry(t *testing.T) {
	cap := &captureSink{}
	ts := sink.NewTransformSink(cap, func(e logpipe.Entry) (logpipe.Entry, bool) {
		return e, false
	})

	e := logpipe.Entry{Level: logpipe.INFO, Message: "drop me", Time: time.Now(), Fields: map[string]interface{}{}}
	if err := ts.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cap.entries) != 0 {
		t.Errorf("expected entry to be dropped, got %d entries", len(cap.entries))
	}
}

func TestTransformSink_PropagatesError(t *testing.T) {
	errSink := &errorSink{err: errors.New("write failed")}
	ts := sink.NewTransformSink(errSink, func(e logpipe.Entry) (logpipe.Entry, bool) {
		return e, true
	})

	e := logpipe.Entry{Level: logpipe.ERROR, Message: "boom", Time: time.Now(), Fields: map[string]interface{}{}}
	if err := ts.Write(e); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestTransformSink_Close(t *testing.T) {
	cap := &captureSink{}
	ts := sink.NewTransformSink(cap, func(e logpipe.Entry) (logpipe.Entry, bool) { return e, true })
	if err := ts.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
	if !cap.closed {
		t.Error("expected underlying sink to be closed")
	}
}
