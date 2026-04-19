package sink

import (
	"errors"

	"github.com/logpipe/logpipe"
)

// CoalesceSink tries each sink in order and returns on the first success.
// If all sinks fail, the errors are joined and returned.
//
// This is useful when you want primary/secondary/tertiary fallback chains
// without the binary limitation of FallbackSink.
type coalesceSink struct {
	sinks []logpipe.Sink
}

// NewCoalesceSink returns a Sink that tries each provided sink in order,
// returning nil on the first successful write. If every sink returns an
// error the combined error is returned to the caller.
func NewCoalesceSink(sinks ...logpipe.Sink) logpipe.Sink {
	return &coalesceSink{sinks: sinks}
}

func (c *coalesceSink) Write(entry logpipe.Entry) error {
	var errs []error
	for _, s := range c.sinks {
		if err := s.Write(entry); err == nil {
			return nil
		} else {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (c *coalesceSink) Close() error {
	var errs []error
	for _, s := range c.sinks {
		if err := s.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
