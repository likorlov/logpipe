package sink

import (
	"errors"

	"github.com/example/logpipe"
)

// MultiSink fans out writes to multiple sinks, collecting all errors.
type MultiSink struct {
	sinks []logpipe.Sink
}

// NewMultiSink creates a sink that writes to all provided sinks.
func NewMultiSink(sinks ...logpipe.Sink) *MultiSink {
	return &MultiSink{sinks: sinks}
}

func (m *MultiSink) Write(e logpipe.Entry) error {
	var errs []error
	for _, s := range m.sinks {
		if err := s.Write(e); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (m *MultiSink) Close() error {
	var errs []error
	for _, s := range m.sinks {
		if err := s.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
