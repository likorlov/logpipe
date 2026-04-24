package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestFlattenSink_FlatEntry(t *testing.T) {
	col := &collectSink{}
	s := sink.NewFlattenSink(col, ".")

	entry := logpipe.Entry{Fields: logpipe.Fields{"level": "info", "msg": "hello"}}
	if err := s.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if col.entries[0].Fields["level"] != "info" {
		t.Errorf("expected level=info, got %v", col.entries[0].Fields["level"])
	}
}

func TestFlattenSink_NestedMap(t *testing.T) {
	col := &collectSink{}
	s := sink.NewFlattenSink(col, ".")

	entry := logpipe.Entry{
		Fields: logpipe.Fields{
			"http": map[string]any{
				"status": 200,
				"method": "GET",
			},
		},
	}
	if err := s.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	f := col.entries[0].Fields
	if f["http.status"] != 200 {
		t.Errorf("expected http.status=200, got %v", f["http.status"])
	}
	if f["http.method"] != "GET" {
		t.Errorf("expected http.method=GET, got %v", f["http.method"])
	}
	if _, ok := f["http"]; ok {
		t.Error("parent key 'http' should have been removed")
	}
}

func TestFlattenSink_DeeplyNested(t *testing.T) {
	col := &collectSink{}
	s := sink.NewFlattenSink(col, "_")

	entry := logpipe.Entry{
		Fields: logpipe.Fields{
			"a": map[string]any{
				"b": map[string]any{
					"c": "deep",
				},
			},
		},
	}
	if err := s.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := col.entries[0].Fields["a_b_c"]; got != "deep" {
		t.Errorf("expected a_b_c=deep, got %v", got)
	}
}

func TestFlattenSink_NoMutationOfOriginal(t *testing.T) {
	col := &collectSink{}
	s := sink.NewFlattenSink(col, ".")

	orig := logpipe.Fields{
		"req": map[string]any{"id": "abc"},
	}
	entry := logpipe.Entry{Fields: orig}
	_ = s.Write(entry)

	if _, ok := orig["req.id"]; ok {
		t.Error("original fields were mutated")
	}
	if _, ok := orig["req"]; !ok {
		t.Error("original 'req' key was removed")
	}
}

func TestFlattenSink_PropagatesError(t *testing.T) {
	expected := errors.New("sink error")
	s := sink.NewFlattenSink(&errSink{err: expected}, ".")
	if err := s.Write(logpipe.Entry{Fields: logpipe.Fields{}}); !errors.Is(err, expected) {
		t.Errorf("expected propagated error, got %v", err)
	}
}

func TestFlattenSink_DefaultSep(t *testing.T) {
	col := &collectSink{}
	s := sink.NewFlattenSink(col, "") // empty sep should default to "."

	entry := logpipe.Entry{
		Fields: logpipe.Fields{
			"x": map[string]any{"y": 1},
		},
	}
	_ = s.Write(entry)
	if col.entries[0].Fields["x.y"] != 1 {
		t.Errorf("expected x.y=1, got %v", col.entries[0].Fields["x.y"])
	}
}

func TestFlattenSink_Close(t *testing.T) {
	col := &collectSink{}
	s := sink.NewFlattenSink(col, ".")
	if err := s.Close(); err != nil {
		t.Errorf("unexpected close error: %v", err)
	}
}
