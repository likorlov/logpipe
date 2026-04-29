package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestSplitterSink_RoutesToMatchingKey(t *testing.T) {
	var gotError, gotDefault logpipe.Entry
	errorSink := &captureSink{writeFn: func(e logpipe.Entry) error { gotError = e; return nil }}
	defaultSink := &captureSink{writeFn: func(e logpipe.Entry) error { gotDefault = e; return nil }}

	s := sink.NewSplitterSink(
		func(e logpipe.Entry) string {
			if lvl, ok := e["level"].(string); ok {
				return lvl
			}
			return "default"
		},
		"default",
		map[string]logpipe.Sink{"error": errorSink, "default": defaultSink},
	)

	if err := s.Write(logpipe.Entry{"level": "error", "msg": "boom"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotError["msg"] != "boom" {
		t.Errorf("error sink did not receive entry")
	}

	if err := s.Write(logpipe.Entry{"level": "info", "msg": "hello"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotDefault["msg"] != "hello" {
		t.Errorf("default sink did not receive info entry")
	}
}

func TestSplitterSink_DropsWhenNoMatchAndNoDefault(t *testing.T) {
	var got logpipe.Entry
	inner := &captureSink{writeFn: func(e logpipe.Entry) error { got = e; return nil }}

	s := sink.NewSplitterSink(
		func(e logpipe.Entry) string { return "unknown" },
		"",
		map[string]logpipe.Sink{"error": inner},
	)

	if err := s.Write(logpipe.Entry{"msg": "dropped"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected entry to be dropped, but sink received it")
	}
}

func TestSplitterSink_PropagatesError(t *testing.T) {
	want := errors.New("write failed")
	inner := &captureSink{writeFn: func(e logpipe.Entry) error { return want }}

	s := sink.NewSplitterSink(
		func(e logpipe.Entry) string { return "x" },
		"",
		map[string]logpipe.Sink{"x": inner},
	)

	if err := s.Write(logpipe.Entry{"msg": "test"}); !errors.Is(err, want) {
		t.Errorf("expected %v, got %v", want, err)
	}
}

func TestSplitterSink_Close(t *testing.T) {
	closed := 0
	mk := func() *captureSink {
		return &captureSink{
			writeFn: func(e logpipe.Entry) error { return nil },
			closeFn: func() error { closed++; return nil },
		}
	}
	s := sink.NewSplitterSink(
		func(e logpipe.Entry) string { return "a" },
		"",
		map[string]logpipe.Sink{"a": mk(), "b": mk()},
	)
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
	if closed != 2 {
		t.Errorf("expected 2 sinks closed, got %d", closed)
	}
}

func TestSplitterSink_PanicOnNilFn(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil fn")
		}
	}()
	sink.NewSplitterSink(nil, "", nil)
}

// captureSink is a minimal test helper used across splitter tests.
type captureSink struct {
	writeFn func(logpipe.Entry) error
	closeFn func() error
}

func (c *captureSink) Write(e logpipe.Entry) error {
	if c.writeFn != nil {
		return c.writeFn(e)
	}
	return nil
}

func (c *captureSink) Close() error {
	if c.closeFn != nil {
		return c.closeFn()
	}
	return nil
}
