package sink_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestTopKSink_TracksFrequency(t *testing.T) {
	var mu sync.Mutex
	var received []logpipe.Entry
	inner := &callbackSink{fn: func(e logpipe.Entry) error {
		mu.Lock()
		received = append(received, e)
		mu.Unlock()
		return nil
	}}

	s := sink.NewTopKSink(inner, "status", 2)

	for _, status := range []string{"ok", "ok", "err", "ok", "err", "warn"} {
		_ = s.Write(logpipe.Entry{Fields: logpipe.Fields{"status": status}})
	}

	top := s.TopK()
	if len(top) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(top))
	}
	if top[0].Value() != "ok" {
		t.Errorf("expected top value 'ok', got %q", top[0].Value())
	}
	if top[1].Value() != "err" {
		t.Errorf("expected second value 'err', got %q", top[1].Value())
	}
	if len(received) != 6 {
		t.Errorf("expected 6 forwarded entries, got %d", len(received))
	}
}

func TestTopKSink_Reset(t *testing.T) {
	inner := &callbackSink{fn: func(e logpipe.Entry) error { return nil }}
	s := sink.NewTopKSink(inner, "msg", 3)

	_ = s.Write(logpipe.Entry{Fields: logpipe.Fields{"msg": "hello"}})
	_ = s.Write(logpipe.Entry{Fields: logpipe.Fields{"msg": "hello"}})
	s.Reset()

	if top := s.TopK(); len(top) != 0 {
		t.Errorf("expected empty after reset, got %d entries", len(top))
	}
}

func TestTopKSink_NonStringFieldIgnored(t *testing.T) {
	inner := &callbackSink{fn: func(e logpipe.Entry) error { return nil }}
	s := sink.NewTopKSink(inner, "code", 5)

	_ = s.Write(logpipe.Entry{Fields: logpipe.Fields{"code": 404}})
	if top := s.TopK(); len(top) != 0 {
		t.Errorf("expected no entries for non-string field, got %d", len(top))
	}
}

func TestTopKSink_PropagatesError(t *testing.T) {
	sentinel := errors.New("inner error")
	inner := &callbackSink{fn: func(e logpipe.Entry) error { return sentinel }}
	s := sink.NewTopKSink(inner, "level", 3)

	err := s.Write(logpipe.Entry{Fields: logpipe.Fields{"level": "info"}})
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestTopKSink_PanicOnZeroK(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for k=0")
		}
	}()
	inner := &callbackSink{fn: func(e logpipe.Entry) error { return nil }}
	sink.NewTopKSink(inner, "field", 0)
}

func TestTopKSink_Close(t *testing.T) {
	closed := false
	inner := &callbackSink{
		fn:      func(e logpipe.Entry) error { return nil },
		closeFn: func() error { closed = true; return nil },
	}
	s := sink.NewTopKSink(inner, "level", 1)
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
	if !closed {
		t.Error("expected inner sink to be closed")
	}
}
