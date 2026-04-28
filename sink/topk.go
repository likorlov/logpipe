package sink

import (
	"sort"
	"sync"

	"github.com/logpipe/logpipe"
)

// topKEntry tracks a field value and its occurrence count.
type topKEntry struct {
	value string
	count int
}

// topKSink tracks the top-K most frequent values for a given field and
// forwards every entry to the inner sink unchanged.
type topKSink struct {
	inner logpipe.Sink
	field string
	k     int
	mu    sync.Mutex
	counts map[string]int
}

// NewTopKSink returns a sink that tracks the top-K most frequent values
// of field across all written entries. Entries are always forwarded to
// inner. Use TopK to retrieve the current leaderboard.
func NewTopKSink(inner logpipe.Sink, field string, k int) *topKSink {
	if k <= 0 {
		panic("logpipe/sink: NewTopKSink k must be > 0")
	}
	return &topKSink{
		inner:  inner,
		field:  field,
		k:      k,
		counts: make(map[string]int),
	}
}

// Write forwards the entry to the inner sink and records the field value.
func (s *topKSink) Write(e logpipe.Entry) error {
	if v, ok := e.Fields[s.field]; ok {
		if str, ok := v.(string); ok && str != "" {
			s.mu.Lock()
			s.counts[str]++
			s.mu.Unlock()
		}
	}
	return s.inner.Write(e)
}

// TopK returns up to k (value, count) pairs ordered by descending count.
func (s *topKSink) TopK() []topKEntry {
	s.mu.Lock()
	entries := make([]topKEntry, 0, len(s.counts))
	for v, c := range s.counts {
		entries = append(entries, topKEntry{value: v, count: c})
	}
	s.mu.Unlock()

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].count != entries[j].count {
			return entries[i].count > entries[j].count
		}
		return entries[i].value < entries[j].value
	})

	if len(entries) > s.k {
		entries = entries[:s.k]
	}
	return entries
}

// Reset clears all accumulated counts.
func (s *topKSink) Reset() {
	s.mu.Lock()
	s.counts = make(map[string]int)
	s.mu.Unlock()
}

// Close closes the inner sink.
func (s *topKSink) Close() error {
	return s.inner.Close()
}
