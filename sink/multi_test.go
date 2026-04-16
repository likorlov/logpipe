package sink_test

import (
	"errors"
	"testing"

	"github.com/example/logpipe"
	"github.com/example/logpipe/sink"
)

type errSink struct{ writeErr error }

func (e *errSink) Write(_ logpipe.Entry) error { return e.writeErr }
func (e *errSink) Close() error               { return nil }

func TestMultiSink_WritesAll(t *testing.T) {
	a, b := &captureSink{}, &captureSink{}
	m := sink.NewMultiSink(a, b)
	defer m.Close()

	_ = m.Write(logpipe.Entry{Message: "broadcast"})
	if a.Len() != 1 || b.Len() != 1 {
		t.Fatalf("expected both sinks to receive entry")
	}
}

func TestMultiSink_CollectsErrors(t *testing.T) {
	e1 := &errSink{writeErr: errors.New("sink1 fail")}
	e2 := &errSink{writeErr: errors.New("sink2 fail")}
	m := sink.NewMultiSink(e1, e2)

	err := m.Write(logpipe.Entry{Message: "x"})
	if err == nil {
		t.Fatal("expected combined error")
	}
	if !errors.Is(err, e1.writeErr) && !errors.Is(err, e2.writeErr) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMultiSink_PartialError(t *testing.T) {
	good := &captureSink{}
	bad := &errSink{writeErr: errors.New("bad")}
	m := sink.NewMultiSink(good, bad)

	err := m.Write(logpipe.Entry{Message: "partial"})
	if err == nil {
		t.Fatal("expected error from bad sink")
	}
	if good.Len() != 1 {
		t.Fatal("good sink should still have received entry")
	}
}
