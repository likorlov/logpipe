package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestReorderSink_FlushOnCapacity(t *testing.T) {
	var got []logpipe.Entry
	inner := &callbackSink{fn: func(e logpipe.Entry) error {
		got = append(got, e)
		return nil
	}}

	s := sink.NewReorderSink(inner, "seq", 3)

	_ = s.Write(logpipe.Entry{Fields: map[string]interface{}{"seq": 3}})
	_ = s.Write(logpipe.Entry{Fields: map[string]interface{}{"seq": 1}})
	// buffer not yet full
	if len(got) != 0 {
		t.Fatalf("expected 0 entries forwarded, got %d", len(got))
	}
	_ = s.Write(logpipe.Entry{Fields: map[string]interface{}{"seq": 2}})
	// now full — should have flushed sorted
	if len(got) != 3 {
		t.Fatalf("expected 3 entries forwarded, got %d", len(got))
	}
	for i, want := range []int{1, 2, 3} {
		if got[i].Fields["seq"] != want {
			t.Errorf("entry[%d] seq = %v, want %d", i, got[i].Fields["seq"], want)
		}
	}
}

func TestReorderSink_FlushOnClose(t *testing.T) {
	var got []logpipe.Entry
	inner := &callbackSink{fn: func(e logpipe.Entry) error {
		got = append(got, e)
		return nil
	}}

	s := sink.NewReorderSink(inner, "seq", 10)
	_ = s.Write(logpipe.Entry{Fields: map[string]interface{}{"seq": 5}})
	_ = s.Write(logpipe.Entry{Fields: map[string]interface{}{"seq": 2}})

	if err := s.Close(); err != nil {
		t.Fatalf("unexpected error on Close: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 entries on close, got %d", len(got))
	}
	if got[0].Fields["seq"] != 2 || got[1].Fields["seq"] != 5 {
		t.Errorf("unexpected order: %v %v", got[0].Fields["seq"], got[1].Fields["seq"])
	}
}

func TestReorderSink_PanicOnZeroCapacity(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for zero capacity")
		}
	}()
	sink.NewReorderSink(&callbackSink{}, "seq", 0)
}

func TestReorderSink_PropagatesError(t *testing.T) {
	wantErr := errors.New("write failed")
	inner := &callbackSink{fn: func(e logpipe.Entry) error { return wantErr }}
	s := sink.NewReorderSink(inner, "seq", 1)
	if err := s.Write(logpipe.Entry{Fields: map[string]interface{}{"seq": 1}}); !errors.Is(err, wantErr) {
		t.Errorf("expected wrapped error, got %v", err)
	}
}

// callbackSink is a test helper that calls fn on each Write.
type callbackSink struct {
	fn func(logpipe.Entry) error
}

func (c *callbackSink) Write(e logpipe.Entry) error {
	if c.fn != nil {
		return c.fn(e)
	}
	return nil
}
func (c *callbackSink) Close() error { return nil }
