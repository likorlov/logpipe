package sink_test

import (
	"errors"
	"testing"

	"github.com/andyh/logpipe"
	"github.com/andyh/logpipe/sink"
)

func TestSchemaSink_PassesValidEntry(t *testing.T) {
	var got []logpipe.Entry
	inner := collectSink(&got)
	s := sink.NewSchemaSink(inner, []sink.SchemaRule{
		{Field: "level", TypeName: "string"},
		{Field: "msg", TypeName: "string"},
	})
	entry := logpipe.Entry{Fields: map[string]any{"level": "info", "msg": "hello"}}
	if err := s.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 entry forwarded, got %d", len(got))
	}
}

func TestSchemaSink_DropsMissingField(t *testing.T) {
	var got []logpipe.Entry
	inner := collectSink(&got)
	s := sink.NewSchemaSink(inner, []sink.SchemaRule{
		{Field: "level"},
		{Field: "msg"},
	})
	entry := logpipe.Entry{Fields: map[string]any{"level": "warn"}} // missing "msg"
	if err := s.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected entry to be dropped, got %d forwarded", len(got))
	}
}

func TestSchemaSink_DropsWrongType(t *testing.T) {
	var got []logpipe.Entry
	inner := collectSink(&got)
	s := sink.NewSchemaSink(inner, []sink.SchemaRule{
		{Field: "count", TypeName: "int"},
	})
	entry := logpipe.Entry{Fields: map[string]any{"count": "not-an-int"}}
	if err := s.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected entry to be dropped due to type mismatch, got %d", len(got))
	}
}

func TestSchemaSink_EmptyRulesPassesAll(t *testing.T) {
	var got []logpipe.Entry
	inner := collectSink(&got)
	s := sink.NewSchemaSink(inner, nil)
	for i := 0; i < 3; i++ {
		_ = s.Write(logpipe.Entry{Fields: map[string]any{"i": i}})
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
}

func TestSchemaSink_Close(t *testing.T) {
	var got []logpipe.Entry
	inner := collectSink(&got)
	s := sink.NewSchemaSink(inner, nil)
	if err := s.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
}

func TestValidateEntry_Valid(t *testing.T) {
	entry := logpipe.Entry{Fields: map[string]any{"level": "info", "ts": 123.4}}
	rules := []sink.SchemaRule{
		{Field: "level", TypeName: "string"},
		{Field: "ts", TypeName: "float64"},
	}
	if err := sink.ValidateEntry(entry, rules); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateEntry_MissingField(t *testing.T) {
	entry := logpipe.Entry{Fields: map[string]any{}}
	rules := []sink.SchemaRule{{Field: "msg"}}
	err := sink.ValidateEntry(entry, rules)
	if !errors.Is(err, sink.ErrSchemaMismatch) {
		t.Fatalf("expected ErrSchemaMismatch, got %v", err)
	}
}

func TestValidateEntry_TypeMismatch(t *testing.T) {
	entry := logpipe.Entry{Fields: map[string]any{"ok": true}}
	rules := []sink.SchemaRule{{Field: "ok", TypeName: "string"}}
	err := sink.ValidateEntry(entry, rules)
	if !errors.Is(err, sink.ErrSchemaMismatch) {
		t.Fatalf("expected ErrSchemaMismatch, got %v", err)
	}
}
