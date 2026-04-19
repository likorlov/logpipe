package sink_test

import (
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestChecksumSink_InjectsField(t *testing.T) {
	col := &collectingSink{}
	s := sink.NewChecksumSink(col, "")
	defer s.Close()

	e := logpipe.Entry{"msg": "hello", "level": "info"}
	if err := s.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
	got := col.entries[0]
	v, ok := got["_checksum"]
	if !ok {
		t.Fatal("expected _checksum field to be present")
	}
	if len(v.(string)) != 16 {
		t.Fatalf("expected 16-char checksum, got %q", v)
	}
}

func TestChecksumSink_CustomField(t *testing.T) {
	col := &collectingSink{}
	s := sink.NewChecksumSink(col, "sig")
	defer s.Close()

	_ = s.Write(logpipe.Entry{"msg": "test"})
	if _, ok := col.entries[0]["sig"]; !ok {
		t.Fatal("expected custom field 'sig'")
	}
}

func TestChecksumSink_NoMutationOfOriginal(t *testing.T) {
	col := &collectingSink{}
	s := sink.NewChecksumSink(col, "")
	defer s.Close()

	orig := logpipe.Entry{"msg": "immutable"}
	_ = s.Write(orig)
	if _, ok := orig["_checksum"]; ok {
		t.Fatal("original entry must not be mutated")
	}
}

func TestChecksumSink_Deterministic(t *testing.T) {
	col := &collectingSink{}
	s := sink.NewChecksumSink(col, "")
	defer s.Close()

	e := logpipe.Entry{"x": "1", "y": "2"}
	_ = s.Write(e)
	_ = s.Write(e)

	a := col.entries[0]["_checksum"].(string)
	b := col.entries[1]["_checksum"].(string)
	if a != b {
		t.Fatalf("expected same checksum for identical entries, got %q vs %q", a, b)
	}
}

func TestChecksumSink_PropagatesError(t *testing.T) {
	s := sink.NewChecksumSink(&errorSink{}, "")
	defer s.Close()
	if err := s.Write(logpipe.Entry{"msg": "fail"}); err == nil {
		t.Fatal("expected error from inner sink")
	}
}
