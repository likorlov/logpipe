package sink_test

import (
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestMaskSink_MasksField(t *testing.T) {
	col := &collectSink{}
	s := sink.NewMaskSink(col, sink.MaskOption{
		Field:      "card",
		KeepPrefix: 4,
		KeepSuffix: 4,
		Mask:       "****",
	})
	_ = s.Write(logpipe.Entry{"card": "1234567890123456"})
	if got := col.entries[0]["card"]; got != "1234****3456" {
		t.Fatalf("expected masked value, got %q", got)
	}
}

func TestMaskSink_DefaultMask(t *testing.T) {
	col := &collectSink{}
	s := sink.NewMaskSink(col, sink.MaskOption{Field: "token", KeepPrefix: 2, KeepSuffix: 0})
	_ = s.Write(logpipe.Entry{"token": "abcdefgh"})
	if got := col.entries[0]["token"]; got != "ab****" {
		t.Fatalf("unexpected value %q", got)
	}
}

func TestMaskSink_ShortValueUnchanged(t *testing.T) {
	col := &collectSink{}
	s := sink.NewMaskSink(col, sink.MaskOption{Field: "pin", KeepPrefix: 2, KeepSuffix: 2})
	_ = s.Write(logpipe.Entry{"pin": "123"})
	if got := col.entries[0]["pin"]; got != "123" {
		t.Fatalf("short value should be unchanged, got %q", got)
	}
}

func TestMaskSink_NoMutationOfOriginal(t *testing.T) {
	col := &collectSink{}
	s := sink.NewMaskSink(col, sink.MaskOption{Field: "secret", KeepPrefix: 1, KeepSuffix: 1})
	orig := logpipe.Entry{"secret": "hello"}
	_ = s.Write(orig)
	if orig["secret"] != "hello" {
		t.Fatal("original entry was mutated")
	}
}

func TestMaskSink_UnrelatedFieldUnchanged(t *testing.T) {
	col := &collectSink{}
	s := sink.NewMaskSink(col, sink.MaskOption{Field: "secret", KeepPrefix: 2, KeepSuffix: 2})
	_ = s.Write(logpipe.Entry{"msg": "hello", "secret": "abcdefgh"})
	if got := col.entries[0]["msg"]; got != "hello" {
		t.Fatalf("unrelated field changed: %q", got)
	}
}

func TestMaskSink_Close(t *testing.T) {
	col := &collectSink{}
	s := sink.NewMaskSink(col)
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
}
