package sink

import (
	"fmt"
	"hash/fnv"
	"sync"

	"github.com/andybar2/logpipe"
)

// SampleRateSink forwards only entries whose key field hashes into the
// configured rate bucket. Unlike SamplingSink (which uses random sampling),
// this sink provides deterministic, consistent sampling based on a field value
// — useful for sampling by user ID, request ID, etc.
type sampleRateSink struct {
	mu      sync.Mutex
	inner   logpipe.Sink
	field   string
	numer   uint32 // entries where hash%denom < numer are forwarded
	denom   uint32
}

// NewSampleRateSink returns a Sink that deterministically forwards a fraction
// (numer/denom) of entries based on the hash of the given field value.
// For example, numer=1, denom=10 forwards ~10% of unique field values.
func NewSampleRateSink(inner logpipe.Sink, field string, numer, denom uint32) (logpipe.Sink, error) {
	if denom == 0 {
		return nil, fmt.Errorf("samplerate: denom must be > 0")
	}
	if numer > denom {
		return nil, fmt.Errorf("samplerate: numer must be <= denom")
	}
	if field == "" {
		field = "id"
	}
	return &sampleRateSink{inner: inner, field: field, numer: numer, denom: denom}, nil
}

func (s *sampleRateSink) Write(entry logpipe.Entry) error {
	val, ok := entry.Fields[s.field]
	if !ok {
		return nil // drop entries missing the key field
	}
	h := fnv.New32a()
	_, _ = fmt.Fprint(h, val)
	if h.Sum32()%s.denom < s.numer {
		s.mu.Lock()
		defer s.mu.Unlock()
		return s.inner.Write(entry)
	}
	return nil
}

func (s *sampleRateSink) Close() error {
	return s.inner.Close()
}
