package sink_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/andygeiss/logpipe"
	"github.com/andygeiss/logpipe/sink"
)

// routeCollect is a thread-safe sink that records written entries.
type routeCollect struct {
	mu      sync.Mutex
	entries []logpipe.Entry
	errOn   bool
}

func (r *routeCollect) Write(e logpipe.Entry) error {
	if r.errOn {
		return errors.New("sink error")
	}
	r.mu.Lock()
	r.entries = append(r.entries, e)
	r.mu.Unlock()
	return nil
}
func (r *routeCollect) Close() error { return nil }
func (r *routeCollect) len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.entries)
}

func TestHashRouteSink_StableRouting(t *testing.T) {
	a, b := &routeCollect{}, &routeCollect{}
	s := sink.NewHashRouteSink("user", a, b)

	e := logpipe.Entry{Fields: map[string]any{"user": "alice"}}
	_ = s.Write(e)
	_ = s.Write(e)
	_ = s.Write(e)

	// All three must land in the same bucket.
	if a.len()+b.len() != 3 {
		t.Fatalf("expected 3 total entries, got %d", a.len()+b.len())
	}
	if a.len() != 3 && b.len() != 3 {
		t.Error("routing is not stable for the same field value")
	}
}

func TestHashRouteSink_DifferentKeysDistribute(t *testing.T) {
	a, b := &routeCollect{}, &routeCollect{}
	s := sink.NewHashRouteSink("user", a, b)

	users := []string{"alice", "bob", "carol", "dave", "eve", "frank"}
	for _, u := range users {
		_ = s.Write(logpipe.Entry{Fields: map[string]any{"user": u}})
	}
	if a.len() == 0 || b.len() == 0 {
		t.Error("expected entries in both sinks; distribution may be broken")
	}
}

func TestHashRouteSink_MissingFieldGoesToFirst(t *testing.T) {
	a, b := &routeCollect{}, &routeCollect{}
	s := sink.NewHashRouteSink("user", a, b)

	_ = s.Write(logpipe.Entry{Fields: map[string]any{"msg": "hello"}})
	if a.len() != 1 {
		t.Errorf("expected entry in first sink, got a=%d b=%d", a.len(), b.len())
	}
}

func TestHashRouteSink_PropagatesError(t *testing.T) {
	bad := &routeCollect{errOn: true}
	s := sink.NewHashRouteSink("user", bad)
	err := s.Write(logpipe.Entry{Fields: map[string]any{"user": "x"}})
	if err == nil {
		t.Error("expected error from inner sink")
	}
}

func TestHashRouteSink_PanicsWithNoSinks(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic with no sinks")
		}
	}()
	sink.NewHashRouteSink("user")
}
