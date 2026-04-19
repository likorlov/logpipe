package sink_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestRoundRobinSink_Distributes(t *testing.T) {
	var mu sync.Mutex
	counts := make([]int, 3)
	makeSink := func(i int) logpipe.Sink {
		return &fnSink{write: func(e logpipe.Entry) error {
			mu.Lock()
			counts[i]++
			mu.Unlock()
			return nil
		}}
	}
	s, err := sink.NewRoundRobinSink(makeSink(0), makeSink(1), makeSink(2))
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 9; i++ {
		if err := s.Write(logpipe.Entry{"msg": "x"}); err != nil {
			t.Fatal(err)
		}
	}
	for i, c := range counts {
		if c != 3 {
			t.Errorf("sink %d: expected 3 writes, got %d", i, c)
		}
	}
}

func TestRoundRobinSink_EmptyError(t *testing.T) {
	_, err := sink.NewRoundRobinSink()
	if err == nil {
		t.Fatal("expected error for empty sinks")
	}
}

func TestRoundRobinSink_PropagatesError(t *testing.T) {
	want := errors.New("boom")
	s, _ := sink.NewRoundRobinSink(&fnSink{write: func(e logpipe.Entry) error { return want }})
	if err := s.Write(logpipe.Entry{}); !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestRoundRobinSink_Close(t *testing.T) {
	closed := 0
	mk := func() logpipe.Sink {
		return &fnSink{write: func(e logpipe.Entry) error { return nil }, close: func() error { closed++; return nil }}
	}
	s, _ := sink.NewRoundRobinSink(mk(), mk(), mk())
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	if closed != 3 {
		t.Fatalf("expected 3 closes, got %d", closed)
	}
}
