package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestAuditSink_InjectsSequence(t *testing.T) {
	var received []logpipe.Entry
	inner := &callbackSink{fn: func(e logpipe.Entry) error { received = append(received, e); return nil }}
	a := sink.NewAuditSink(inner, "")

	for i := 0; i < 3; i++ {
		_ = a.Write(logpipe.Entry{Level: logpipe.LevelInfo, Message: "msg", Fields: map[string]any{}})
	}

	if len(received) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(received))
	}
	for i, e := range received {
		want := fmt.Sprintf("%d", i+1)
		if e.Fields["_audit_seq"] != want {
			t.Errorf("entry %d: want _audit_seq=%q got %q", i, want, e.Fields["_audit_seq"])
		}
	}
}

func TestAuditSink_CustomField(t *testing.T) {
	inner := &callbackSink{fn: func(e logpipe.Entry) error { return nil }}
	a := sink.NewAuditSink(inner, "seq_no")
	_ = a.Write(logpipe.Entry{Level: logpipe.LevelInfo, Message: "x", Fields: map[string]any{}})
	entries := a.Entries()
	if _, ok := entries[0].Fields["seq_no"]; !ok {
		t.Fatal("expected seq_no field")
	}
}

func TestAuditSink_NoMutationOfOriginal(t *testing.T) {
	inner := &callbackSink{fn: func(e logpipe.Entry) error { return nil }}
	a := sink.NewAuditSink(inner, "")
	orig := logpipe.Entry{Level: logpipe.LevelInfo, Message: "m", Fields: map[string]any{"k": "v"}}
	_ = a.Write(orig)
	if _, ok := orig.Fields["_audit_seq"]; ok {
		t.Fatal("original entry was mutated")
	}
}

func TestAuditSink_Reset(t *testing.T) {
	inner := &callbackSink{fn: func(e logpipe.Entry) error { return nil }}
	a := sink.NewAuditSink(inner, "")
	_ = a.Write(logpipe.Entry{Level: logpipe.LevelInfo, Message: "m", Fields: map[string]any{}})
	a.Reset()
	if len(a.Entries()) != 0 {
		t.Fatal("expected empty entries after reset")
	}
	_ = a.Write(logpipe.Entry{Level: logpipe.LevelInfo, Message: "m2", Fields: map[string]any{}})
	if a.Entries()[0].Fields["_audit_seq"] != "1" {
		t.Fatal("sequence did not reset to 1")
	}
}

func TestAuditSink_PropagatesError(t *testing.T) {
	expected := errors.New("boom")
	inner := &callbackSink{fn: func(e logpipe.Entry) error { return expected }}
	a := sink.NewAuditSink(inner, "")
	if err := a.Write(logpipe.Entry{Level: logpipe.LevelInfo, Message: "m", Fields: map[string]any{}}); err != expected {
		t.Fatalf("expected propagated error, got %v", err)
	}
}
