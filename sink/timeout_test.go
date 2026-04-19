package sink_test

import (
	"errors"
	"testing"
	"time"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

type slowSink struct {
	delay time.Duration
	err   error
}

func (s *slowSink) Write(_ logpipe.Entry) error {
	time.Sleep(s.delay)
	return s.err
}
func (s *slowSink) Close() error { return nil }

func TestTimeoutSink_PassesOnTime(t *testing.T) {
	inner := &slowSink{delay: 10 * time.Millisecond}
	s := sink.NewTimeoutSink(inner, 200*time.Millisecond)
	defer s.Close()

	if err := s.Write(logpipe.Entry{Message: "fast"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTimeoutSink_TimesOut(t *testing.T) {
	inner := &slowSink{delay: 300 * time.Millisecond}
	s := sink.NewTimeoutSink(inner, 50*time.Millisecond)
	defer s.Close()

	err := s.Write(logpipe.Entry{Message: "slow"})
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestTimeoutSink_PropagatesInnerError(t *testing.T) {
	want := errors.New("inner failure")
	inner := &slowSink{delay: 0, err: want}
	s := sink.NewTimeoutSink(inner, 200*time.Millisecond)
	defer s.Close()

	if err := s.Write(logpipe.Entry{Message: "err"}); !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestTimeoutSink_Close(t *testing.T) {
	inner := &slowSink{}
	s := sink.NewTimeoutSink(inner, time.Second)
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
}
