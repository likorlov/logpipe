package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestCeilingSink_AllowsUnderCeiling(t *testing.T) {
	col := &collectSink{}
	s := sink.NewCeilingSink(col, "count", 10)

	if err := s.Write(logpipe.Entry{"count": 5}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
}

func TestCeilingSink_AllowsAtCeiling(t *testing.T) {
	col := &collectSink{}
	s := sink.NewCeilingSink(col, "count", 10)

	if err := s.Write(logpipe.Entry{"count": 10}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
}

func TestCeilingSink_DropsAboveCeiling(t *testing.T) {
	col := &collectSink{}
	s := sink.NewCeilingSink(col, "count", 10)

	if err := s.Write(logpipe.Entry{"count": 11}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(col.entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(col.entries))
	}
}

func TestCeilingSink_MissingFieldDrops(t *testing.T) {
	col := &collectSink{}
	s := sink.NewCeilingSink(col, "count", 10)

	if err := s.Write(logpipe.Entry{"msg": "hello"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(col.entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(col.entries))
	}
}

func TestCeilingSink_NonNumericFieldDrops(t *testing.T) {
	col := &collectSink{}
	s := sink.NewCeilingSink(col, "count", 10)

	if err := s.Write(logpipe.Entry{"count": "high"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(col.entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(col.entries))
	}
}

func TestCeilingSink_DefaultField(t *testing.T) {
	col := &collectSink{}
	s := sink.NewCeilingSink(col, "", 5)

	if err := s.Write(logpipe.Entry{"value": float64(3)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(col.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(col.entries))
	}
}

func TestCeilingSink_PropagatesError(t *testing.T) {
	sentinel := errors.New("inner error")
	errSink := &errorSink{err: sentinel}
	s := sink.NewCeilingSink(errSink, "count", 100)

	err := s.Write(logpipe.Entry{"count": 1})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestCeilingSink_PanicOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner")
		}
	}()
	sink.NewCeilingSink(nil, "count", 10)
}

func TestCeilingSink_PanicOnNegativeCeiling(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative ceiling")
		}
	}()
	sink.NewCeilingSink(&collectSink{}, "count", -1)
}

func TestCeilingSink_Close(t *testing.T) {
	col := &collectSink{}
	s := sink.NewCeilingSink(col, "count", 10)
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
}
