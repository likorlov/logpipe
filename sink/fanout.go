package sink

import (
	"fmt"
	"strings"

	"github.com/logpipe/logpipe"
)

// fanoutSink writes each entry to all sinks concurrently and collects errors.
type fanoutSink struct {
	sinks []logpipe.Sink
}

// NewFanoutSink returns a Sink that writes each log entry to all provided sinks
// concurrently. Unlike MultiSink, FanoutSink dispatches writes in parallel and
// waits for all of them to complete before returning. All errors are collected
// and returned as a single combined error.
//
// Example:
//
//	sink.NewFanoutSink(
//		sink.NewConsoleSink(),
//		sink.NewFileSink("/var/log/app.log"),
//	)
func NewFanoutSink(sinks ...logpipe.Sink) logpipe.Sink {
	return &fanoutSink{sinks: sinks}
}

func (f *fanoutSink) Write(entry logpipe.Entry) error {
	type result struct {
		idx int
		err error
	}

	ch := make(chan result, len(f.sinks))
	for i, s := range f.sinks {
		go func(idx int, s logpipe.Sink) {
			ch <- result{idx: idx, err: s.Write(entry)}
		}(i, s)
	}

	var errs []string
	for range f.sinks {
		if r := <-ch; r.err != nil {
			errs = append(errs, fmt.Sprintf("sink[%d]: %v", r.idx, r.err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("fanout errors: %s", strings.Join(errs, "; "))
	}
	return nil
}

func (f *fanoutSink) Close() error {
	var errs []string
	for i, s := range f.sinks {
		if err := s.Close(); err != nil {
			errs = append(errs, fmt.Sprintf("sink[%d]: %v", i, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("fanout close errors: %s", strings.Join(errs, "; "))
	}
	return nil
}
