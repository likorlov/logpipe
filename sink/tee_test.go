package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestTeeSink_BothReceiveEntry(t *testing.T) {
	var a, b []logpipe.Entry
	pa := sink.NewTransformSink(collectSink(&a), func(e logpipe.Entry) (logpipe.Entry, bool) { return e, true })
	pb := sink.NewTransformSink(collectSink(&b), func(e logpipe.Entry) (logpipe.Entry, bool) { return e, true })
	tee := sink.NewTeeSink(pa, pb)

	entry := logpipe.Entry{Message: "hello", Level: logpipe.INFO}
	if err := tee.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(a) != 1 || len(b) != 1 {
		t.Fatalf("expected both sinks to receive entry, got a=%d b=%d", len(a), len(b))
	}
}

func TestTeeSink_PrimaryErrorSkipsSecondary(t *testing.T) {
	var b []logpipe.Entry
	primErr := errors.New("primary fail")
	primary := &errSink{err: primErr}
	secondary := collectSink(&b)
	tee := sink.NewTeeSink(primary, secondary)

	err := tee.Write(logpipe.Entry{Message: "x", Level: logpipe.WARN})
	if !errors.Is(err, primErr) {
		t.Fatalf("expected primary error, got %v", err)
	}
	if len(b) != 0 {
		t.Fatal("secondary should not receive entry when primary fails")
	}
}

func TestTeeSink_SecondaryErrorReturned(t *testing.T) {
	var a []logpipe.Entry
	secErr := errors.New("secondary fail")
	primary := collectSink(&a)
	secondary := &errSink{err: secErr}
	tee := sink.NewTeeSink(primary, secondary)

	err := tee.Write(logpipe.Entry{Message: "y", Level: logpipe.ERROR})
	if !errors.Is(err, secErr) {
		t.Fatalf("expected secondary error, got %v", err)
	}
	if len(a) != 1 {
		t.Fatal("primary should still have received the entry")
	}
}

func TestTeeSink_Close(t *testing.T) {
	var a, b []logpipe.Entry
	tee := sink.NewTeeSink(collectSink(&a), collectSink(&b))
	if err := tee.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
}

// helpers shared across sink tests

type collectSinkImpl struct{ entries *[]logpipe.Entry }

func collectSink(out *[]logpipe.Entry) logpipe.Sink { return &collectSinkImpl{entries: out} }
func (c *collectSinkImpl) Write(e logpipe.Entry) error { *c.entries = append(*c.entries, e); return nil }
func (c *collectSinkImpl) Close() error               { return nil }

type errSink struct{ err error }

func (e *errSink) Write(_ logpipe.Entry) error { return e.err }
func (e *errSink) Close() error               { return nil }
