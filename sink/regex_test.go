package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestRegexSink_MatchForwards(t *testing.T) {
	col := &collectSink{}
	s, err := sink.NewRegexSink(col, "msg", `^error`, false)
	if err != nil {
		t.Fatal(err)
	}
	e := logpipe.Entry{Fields: map[string]any{"msg": "error: something went wrong"}}
	if err := s.Write(e); err != nil {
		t.Fatal(err)
	}
	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
}

func TestRegexSink_NoMatchDrops(t *testing.T) {
	col := &collectSink{}
	s, _ := sink.NewRegexSink(col, "msg", `^error`, false)
	e := logpipe.Entry{Fields: map[string]any{"msg": "info: all good"}}
	if err := s.Write(e); err != nil {
		t.Fatal(err)
	}
	if len(col.entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(col.entries))
	}
}

func TestRegexSink_InvertDropsMatch(t *testing.T) {
	col := &collectSink{}
	s, _ := sink.NewRegexSink(col, "msg", `^error`, true)
	e := logpipe.Entry{Fields: map[string]any{"msg": "error: oops"}}
	if err := s.Write(e); err != nil {
		t.Fatal(err)
	}
	if len(col.entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(col.entries))
	}
}

func TestRegexSink_MissingFieldForwards(t *testing.T) {
	col := &collectSink{}
	s, _ := sink.NewRegexSink(col, "msg", `^error`, false)
	e := logpipe.Entry{Fields: map[string]any{"other": "value"}}
	if err := s.Write(e); err != nil {
		t.Fatal(err)
	}
	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
}

func TestRegexSink_NonStringFieldForwards(t *testing.T) {
	col := &collectSink{}
	s, _ := sink.NewRegexSink(col, "code", `^5`, false)
	e := logpipe.Entry{Fields: map[string]any{"code": 500}}
	if err := s.Write(e); err != nil {
		t.Fatal(err)
	}
	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
}

func TestRegexSink_InvalidPatternError(t *testing.T) {
	col := &collectSink{}
	_, err := sink.NewRegexSink(col, "msg", `[invalid`, false)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestRegexSink_PropagatesError(t *testing.T) {
	want := errors.New("inner error")
	fail := &failSink{err: want}
	s, _ := sink.NewRegexSink(fail, "msg", `.*`, false)
	e := logpipe.Entry{Fields: map[string]any{"msg": "hello"}}
	if err := s.Write(e); !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestRegexSink_Close(t *testing.T) {
	col := &collectSink{}
	s, _ := sink.NewRegexSink(col, "msg", `.*`, false)
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	if !col.closed {
		t.Fatal("expected inner sink to be closed")
	}
}
