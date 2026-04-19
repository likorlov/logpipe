package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestPrioritySink_RoutesToFirstMatch(t *testing.T) {
	var gotDebug, gotError []logpipe.Entry

	debugSink := &collectSink{fn: func(e logpipe.Entry) error { gotDebug = append(gotDebug, e); return nil }}
	errorSink := &collectSink{fn: func(e logpipe.Entry) error { gotError = append(gotError, e); return nil }}

	p := sink.NewPrioritySink()
	p.Add(logpipe.LevelError, errorSink)
	p.Add(logpipe.LevelDebug, debugSink)

	p.Write(logpipe.Entry{Level: logpipe.LevelError, Message: "boom"})
	p.Write(logpipe.Entry{Level: logpipe.LevelInfo, Message: "info"})
	p.Write(logpipe.Entry{Level: logpipe.LevelDebug, Message: "dbg"})

	if len(gotError) != 1 || gotError[0].Message != "boom" {
		t.Fatalf("expected error sink to receive 1 entry, got %v", gotError)
	}
	if len(gotDebug) != 2 {
		t.Fatalf("expected debug sink to receive 2 entries, got %v", gotDebug)
	}
}

func TestPrioritySink_NoMatchDrops(t *testing.T) {
	var received []logpipe.Entry
	s := &collectSink{fn: func(e logpipe.Entry) error { received = append(received, e); return nil }}

	p := sink.NewPrioritySink()
	p.Add(logpipe.LevelError, s)

	p.Write(logpipe.Entry{Level: logpipe.LevelDebug, Message: "quiet"})

	if len(received) != 0 {
		t.Fatalf("expected no entries, got %v", received)
	}
}

func TestPrioritySink_PropagatesError(t *testing.T) {
	want := errors.New("sink down")
	s := &collectSink{fn: func(e logpipe.Entry) error { return want }}

	p := sink.NewPrioritySink()
	p.Add(logpipe.LevelDebug, s)

	if err := p.Write(logpipe.Entry{Level: logpipe.LevelInfo}); !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestPrioritySink_Close(t *testing.T) {
	p := sink.NewPrioritySink()
	p.Add(logpipe.LevelDebug, &collectSink{fn: func(e logpipe.Entry) error { return nil }})
	if err := p.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
