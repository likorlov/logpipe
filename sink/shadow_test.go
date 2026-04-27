package sink_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

// shadowCollector is a simple in-memory sink used by shadow tests.
type shadowCollector struct {
	mu      sync.Mutex
	entries []logpipe.Entry
	err     error
}

func (c *shadowCollector) Write(e logpipe.Entry) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err != nil {
		return c.err
	}
	copy := make(logpipe.Entry, len(e))
	for k, v := range e {
		copy[k] = v
	}
	c.entries = append(c.entries, copy)
	return nil
}
func (c *shadowCollector) Close() error { return nil }
func (c *shadowCollector) all() []logpipe.Entry {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.entries
}

func TestShadowSink_BothReceiveEntry(t *testing.T) {
	primary := &shadowCollector{}
	shadow := &shadowCollector{}
	s := sink.NewShadowSink(primary, shadow)

	entry := logpipe.Entry{"msg": "hello"}
	if err := s.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(primary.all()) != 1 {
		t.Errorf("primary: expected 1 entry, got %d", len(primary.all()))
	}
	if len(shadow.all()) != 1 {
		t.Errorf("shadow: expected 1 entry, got %d", len(shadow.all()))
	}
}

func TestShadowSink_ShadowErrorIgnored(t *testing.T) {
	primary := &shadowCollector{}
	shadow := &shadowCollector{err: errors.New("shadow boom")}
	s := sink.NewShadowSink(primary, shadow)

	if err := s.Write(logpipe.Entry{"msg": "ok"}); err != nil {
		t.Fatalf("shadow error must not propagate, got: %v", err)
	}
	if len(primary.all()) != 1 {
		t.Errorf("expected primary to receive entry")
	}
}

func TestShadowSink_PrimaryErrorReturned(t *testing.T) {
	expected := errors.New("primary boom")
	primary := &shadowCollector{err: expected}
	shadow := &shadowCollector{}
	s := sink.NewShadowSink(primary, shadow)

	if err := s.Write(logpipe.Entry{"msg": "x"}); !errors.Is(err, expected) {
		t.Errorf("expected primary error, got %v", err)
	}
}

func TestShadowSink_NoMutationOfShadowEntry(t *testing.T) {
	primary := &shadowCollector{}
	shadow := &shadowCollector{}
	s := sink.NewShadowSink(primary, shadow)

	original := logpipe.Entry{"key": "value"}
	_ = s.Write(original)

	// Mutate what primary received; shadow copy must be unaffected.
	primary.all()[0]["key"] = "mutated"
	if shadow.all()[0]["key"] != "value" {
		t.Errorf("shadow entry was mutated via primary")
	}
}

func TestShadowSink_Close(t *testing.T) {
	primary := &shadowCollector{}
	shadow := &shadowCollector{}
	s := sink.NewShadowSink(primary, shadow)
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
}
