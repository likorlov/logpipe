package sink

import (
	"fmt"
	"hash/fnv"

	"github.com/andygeiss/logpipe"
)

// HashRouteSink routes each log entry to one of N sinks based on a stable
// hash of a chosen field value. Entries missing the field go to sink 0.
type hashRouteSink struct {
	field string
	sinks []logpipe.Sink
}

// NewHashRouteSink returns a Sink that shards entries across the provided
// sinks by hashing the value of field. At least one sink must be supplied.
func NewHashRouteSink(field string, sinks ...logpipe.Sink) logpipe.Sink {
	if len(sinks) == 0 {
		panic("hashroute: at least one sink required")
	}
	return &hashRouteSink{field: field, sinks: sinks}
}

func (h *hashRouteSink) Write(e logpipe.Entry) error {
	idx := 0
	if v, ok := e.Fields[h.field]; ok {
		hash := fnv.New32a()
		_, _ = hash.Write([]byte(fmt.Sprintf("%v", v)))
		idx = int(hash.Sum32()) % len(h.sinks)
	}
	return h.sinks[idx].Write(e)
}

func (h *hashRouteSink) Close() error {
	var first error
	for _, s := range h.sinks {
		if err := s.Close(); err != nil && first == nil {
			first = err
		}
	}
	return first
}
