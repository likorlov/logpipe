package sink

import (
	"sync"
	"time"

	"github.com/logpipe/logpipe"
)

// cacheSink wraps an inner sink and caches the last write result for identical
// entries within a TTL window. Duplicate entries that hit the cache are
// forwarded to the inner sink but reuse the previously computed error result,
// avoiding redundant work for high-frequency identical log lines.
type cacheSink struct {
	inner logpipe.Sink
	ttl   time.Duration
	mu    sync.Mutex
	cache map[string]cacheEntry
}

type cacheEntry struct {
	err     error
	expires time.Time
}

// NewCacheSink returns a Sink that caches the result of writing an entry for
// the given TTL. Entries are keyed by their message field. If the same message
// is received within the TTL the cached error (nil on success) is returned
// without forwarding to the inner sink.
//
//	cached := sink.NewCacheSink(inner, 5*time.Second)
func NewCacheSink(inner logpipe.Sink, ttl time.Duration) logpipe.Sink {
	return &cacheSink{
		inner: inner,
		ttl:   ttl,
		cache: make(map[string]cacheEntry),
	}
}

func (s *cacheSink) Write(entry logpipe.Entry) error {
	key, _ := entry.Fields["message"].(string)
	if key == "" {
		return s.inner.Write(entry)
	}

	now := time.Now()

	s.mu.Lock()
	if ce, ok := s.cache[key]; ok && now.Before(ce.expires) {
		s.mu.Unlock()
		return ce.err
	}
	s.mu.Unlock()

	err := s.inner.Write(entry)

	s.mu.Lock()
	s.cache[key] = cacheEntry{err: err, expires: now.Add(s.ttl)}
	s.mu.Unlock()

	return err
}

func (s *cacheSink) Close() error {
	s.mu.Lock()
	s.cache = make(map[string]cacheEntry)
	s.mu.Unlock()
	return s.inner.Close()
}
