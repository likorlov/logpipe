package sink_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestNormalizeSink_LowercasesKeys(t *testing.T) {
	var got logpipe.Entry
	collect := &captureSink{writeFn: func(e logpipe.Entry) error { got = e; return nil }}
	s := sink.NewNormalizeSink(collect, nil)

	_ = s.Write(logpipe.Entry{
		Level:   logpipe.LevelInfo,
		Message: "hello",
		Fields:  map[string]any{"UserID": 1, "RequestID": "abc"},
	})

	if _, ok := got.Fields["userid"]; !ok {
		t.Error("expected key 'userid'")
	}
	if _, ok := got.Fields["requestid"]; !ok {
		t.Error("expected key 'requestid'")
	}
}

func TestNormalizeSink_CustomFn(t *testing.T) {
	var got logpipe.Entry
	collect := &captureSink{writeFn: func(e logpipe.Entry) error { got = e; return nil }}
	s := sink.NewNormalizeSink(collect, strings.ToUpper)

	_ = s.Write(logpipe.Entry{
		Level:   logpipe.LevelInfo,
		Message: "msg",
		Fields:  map[string]any{"foo": "bar"},
	})

	if _, ok := got.Fields["FOO"]; !ok {
		t.Error("expected key 'FOO'")
	}
}

func TestNormalizeSink_NoMutationOfOriginal(t *testing.T) {
	collect := &captureSink{writeFn: func(e logpipe.Entry) error { return nil }}
	s := sink.NewNormalizeSink(collect, nil)

	orig := logpipe.Entry{
		Level:   logpipe.LevelInfo,
		Message: "msg",
		Fields:  map[string]any{"MyKey": 42},
	}
	_ = s.Write(orig)

	if _, ok := orig.Fields["MyKey"]; !ok {
		t.Error("original entry was mutated")
	}
}

func TestNormalizeSink_PropagatesError(t *testing.T) {
	expected := errors.New("inner error")
	collect := &captureSink{writeFn: func(e logpipe.Entry) error { return expected }}
	s := sink.NewNormalizeSink(collect, nil)

	err := s.Write(logpipe.Entry{Fields: map[string]any{}})
	if !errors.Is(err, expected) {
		t.Errorf("expected %v, got %v", expected, err)
	}
}

func TestNormalizeSink_Close(t *testing.T) {
	closed := false
	collect := &captureSink{closeFn: func() error { closed = true; return nil }}
	s := sink.NewNormalizeSink(collect, nil)
	_ = s.Close()
	if !closed {
		t.Error("expected inner sink to be closed")
	}
}
