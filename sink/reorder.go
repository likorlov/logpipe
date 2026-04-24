package sink

import (
	"sort"
	"sync"

	"github.com/logpipe/logpipe"
)

// reorderSink buffers up to capacity entries and flushes them sorted by a
// numeric field (ascending). Entries are flushed when the buffer is full or
// when Close is called.
type reorderSink struct {
	inner    logpipe.Sink
	field    string
	capacity int
	mu       sync.Mutex
	buf      []logpipe.Entry
}

// NewReorderSink returns a Sink that collects up to capacity entries and
// forwards them to inner sorted by the given numeric field (ascending).
// When the buffer reaches capacity it is flushed immediately. Any remaining
// buffered entries are flushed on Close.
func NewReorderSink(inner logpipe.Sink, field string, capacity int) logpipe.Sink {
	if capacity <= 0 {
		panic("logpipe/sink: NewReorderSink capacity must be > 0")
	}
	return &reorderSink{
		inner:    inner,
		field:    field,
		capacity: capacity,
		buf:      make([]logpipe.Entry, 0, capacity),
	}
}

func (s *reorderSink) Write(e logpipe.Entry) error {
	s.mu.Lock()
	s.buf = append(s.buf, e)
	if len(s.buf) >= s.capacity {
		err := s.flush()
		s.mu.Unlock()
		return err
	}
	s.mu.Unlock()
	return nil
}

func (s *reorderSink) flush() error {
	sort.SliceStable(s.buf, func(i, j int) bool {
		return numVal(s.buf[i].Fields[s.field]) < numVal(s.buf[j].Fields[s.field])
	})
	for _, e := range s.buf {
		if err := s.inner.Write(e); err != nil {
			s.buf = s.buf[:0]
			return err
		}
	}
	s.buf = s.buf[:0]
	return nil
}

func (s *reorderSink) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.buf) > 0 {
		if err := s.flush(); err != nil {
			return err
		}
	}
	return s.inner.Close()
}

// numVal coerces a field value to float64 for comparison purposes.
func numVal(v interface{}) float64 {
	switch n := v.(type) {
	case int:
		return float64(n)
	case int64:
		return float64(n)
	case float64:
		return n
	case float32:
		return float64(n)
	}
	return 0
}
