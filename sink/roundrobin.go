package sink

import (
	"fmt"
	"sync/atomic"

	"github.com/logpipe/logpipe"
)

// RoundRobinSink distributes log entries across multiple sinks in round-robin order.
type roundRobinSink struct {
	sinks   []logpipe.Sink
	counter atomic.Uint64
}

// NewRoundRobinSink returns a Sink that forwards each entry to the next sink
// in rotation. At least one sink must be provided.
func NewRoundRobinSink(sinks ...logpipe.Sink) (logpipe.Sink, error) {
	if len(sinks) == 0 {
		return nil, fmt.Errorf("logpipe/sink: NewRoundRobinSink requires at least one sink")
	}
	return &roundRobinSink{sinks: sinks}, nil
}

func (r *roundRobinSink) Write(entry logpipe.Entry) error {
	n := r.counter.Add(1) - 1
	idx := n % uint64(len(r.sinks))
	return r.sinks[idx].Write(entry)
}

func (r *roundRobinSink) Close() error {
	var first error
	for _, s := range r.sinks {
		if err := s.Close(); err != nil && first == nil {
			first = err
		}
	}
	return first
}
