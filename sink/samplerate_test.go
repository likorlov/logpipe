package sink_test

import (
	"fmt"
	"testing"

	"github.com/andybar2/logpipe"
	"github.com/andybar2/logpipe/sink"
)

func TestSampleRateSink_InvalidDenom(t *testing.T) {
	_, err := sink.NewSampleRateSink(&captureSink{}, "id", 1, 0)
	if err == nil {
		t.Fatal("expected error for denom=0")
	}
}

func TestSampleRateSink_NumerGTDenom(t *testing.T) {
	_, err := sink.NewSampleRateSink(&captureSink{}, "id", 5, 3)
	if err == nil {
		t.Fatal("expected error for numer > denom")
	}
}

func TestSampleRateSink_AllPass(t *testing.T) {
	cap := &captureSink{}
	s, err := sink.NewSampleRateSink(cap, "id", 10, 10)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 20; i++ {
		_ = s.Write(logpipe.Entry{Fields: map[string]any{"id": fmt.Sprintf("user-%d", i)}})
	}
	if len(cap.entries) != 20 {
		t.Fatalf("expected 20 entries, got %d", len(cap.entries))
	}
}

func TestSampleRateSink_NonePass(t *testing.T) {
	cap := &captureSink{}
	s, err := sink.NewSampleRateSink(cap, "id", 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 20; i++ {
		_ = s.Write(logpipe.Entry{Fields: map[string]any{"id": fmt.Sprintf("user-%d", i)}})
	}
	if len(cap.entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(cap.entries))
	}
}

func TestSampleRateSink_MissingFieldDrops(t *testing.T) {
	cap := &captureSink{}
	s, _ := sink.NewSampleRateSink(cap, "id", 10, 10)
	_ = s.Write(logpipe.Entry{Fields: map[string]any{"msg": "hello"}})
	if len(cap.entries) != 0 {
		t.Fatal("expected entry to be dropped when field missing")
	}
}

func TestSampleRateSink_Deterministic(t *testing.T) {
	cap1 := &captureSink{}
	cap2 := &captureSink{}
	s1, _ := sink.NewSampleRateSink(cap1, "id", 3, 10)
	s2, _ := sink.NewSampleRateSink(cap2, "id", 3, 10)
	for i := 0; i < 100; i++ {
		e := logpipe.Entry{Fields: map[string]any{"id": fmt.Sprintf("key-%d", i)}}
		_ = s1.Write(e)
		_ = s2.Write(e)
	}
	if len(cap1.entries) != len(cap2.entries) {
		t.Fatalf("expected deterministic results: %d vs %d", len(cap1.entries), len(cap2.entries))
	}
}
