package sink_test

import (
	"testing"
	"time"

	"github.com/your/logpipe"
	"github.com/your/logpipe/sink"
)

func TestRateLimitSink_AllowsUnderLimit(t *testing.T) {
	rec := &recorder{}
	rl := sink.NewRateLimitSink(rec, 5, time.Second)
	for i := 0; i < 5; i++ {
		if err := rl.Write(logpipe.Entry{Message: "msg"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if len(rec.entries) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(rec.entries))
	}
}

func TestRateLimitSink_DropsOverLimit(t *testing.T) {
	rec := &recorder{}
	rl := sink.NewRateLimitSink(rec, 3, time.Second)
	for i := 0; i < 6; i++ {
		_ = rl.Write(logpipe.Entry{Message: "msg"})
	}
	if len(rec.entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(rec.entries))
	}
}

func TestRateLimitSink_ResetsAfterInterval(t *testing.T) {
	rec := &recorder{}
	rl := sink.NewRateLimitSink(rec, 2, 50*time.Millisecond)
	for i := 0; i < 2; i++ {
		_ = rl.Write(logpipe.Entry{Message: "msg"})
	}
	time.Sleep(60 * time.Millisecond)
	for i := 0; i < 2; i++ {
		_ = rl.Write(logpipe.Entry{Message: "msg"})
	}
	if len(rec.entries) != 4 {
		t.Fatalf("expected 4 entries after reset, got %d", len(rec.entries))
	}
}

func TestRateLimitSink_Close(t *testing.T) {
	rec := &recorder{}
	rl := sink.NewRateLimitSink(rec, 10, time.Second)
	if err := rl.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
	if !rec.closed {
		t.Fatal("expected wrapped sink to be closed")
	}
}
