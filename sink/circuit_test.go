package sink_test

import (
	"errors"
	"testing"
	"time"

	"github.com/your/logpipe"
	"github.com/your/logpipe/sink"
)

type failSink struct {
	calls  int
	failOn int // fail when calls <= failOn
}

func (f *failSink) Write(_ logpipe.Entry) error {
	f.calls++
	if f.calls <= f.failOn {
		return errors.New("downstream error")
	}
	return nil
}
func (f *failSink) Close() error { return nil }

func TestCircuitSink_ClosedOnSuccess(t *testing.T) {
	inner := &failSink{failOn: 0}
	s := sink.NewCircuitSink(inner, 3, time.Minute)
	entry := logpipe.Entry{Message: "ok"}

	if err := s.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCircuitSink_OpensAfterMaxFails(t *testing.T) {
	inner := &failSink{failOn: 10}
	s := sink.NewCircuitSink(inner, 3, time.Minute)
	entry := logpipe.Entry{Message: "x"}

	for i := 0; i < 3; i++ {
		s.Write(entry) //nolint
	}

	err := s.Write(entry)
	if !errors.Is(err, sink.ErrCircuitOpen) {
		t.Fatalf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitSink_ResetsAfterCooldown(t *testing.T) {
	inner := &failSink{failOn: 3} // first 3 calls fail, then succeed
	s := sink.NewCircuitSink(inner, 3, 10*time.Millisecond)
	entry := logpipe.Entry{Message: "x"}

	for i := 0; i < 3; i++ {
		s.Write(entry) //nolint
	}

	// circuit open
	if err := s.Write(entry); !errors.Is(err, sink.ErrCircuitOpen) {
		t.Fatal("expected open circuit")
	}

	time.Sleep(20 * time.Millisecond)

	// after cooldown, next write should reach inner (call 4 succeeds)
	if err := s.Write(entry); err != nil {
		t.Fatalf("expected circuit reset, got %v", err)
	}
}

func TestCircuitSink_Close(t *testing.T) {
	inner := &failSink{}
	s := sink.NewCircuitSink(inner, 3, time.Minute)
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
}
