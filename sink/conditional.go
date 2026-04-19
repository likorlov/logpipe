package sink

import "github.com/logpipe/logpipe"

// ConditionalSink forwards log entries to one of two sinks based on a predicate.
// If the predicate returns true the entry goes to the "then" sink, otherwise
// to the "else" sink. Either branch may be nil, in which case matching entries
// are silently dropped.
type ConditionalSink struct {
	predicate func(logpipe.Entry) bool
	thenSink  logpipe.Sink
	elseSink  logpipe.Sink
}

// NewConditionalSink creates a ConditionalSink.
// thenSink receives entries where predicate returns true.
// elseSink receives entries where predicate returns false.
// Either sink may be nil.
func NewConditionalSink(predicate func(logpipe.Entry) bool, thenSink, elseSink logpipe.Sink) *ConditionalSink {
	if predicate == nil {
		panic("logpipe/sink: ConditionalSink predicate must not be nil")
	}
	return &ConditionalSink{predicate: predicate, thenSink: thenSink, elseSink: elseSink}
}

// Write routes the entry to the appropriate branch sink.
func (c *ConditionalSink) Write(entry logpipe.Entry) error {
	if c.predicate(entry) {
		if c.thenSink != nil {
			return c.thenSink.Write(entry)
		}
		return nil
	}
	if c.elseSink != nil {
		return c.elseSink.Write(entry)
	}
	return nil
}

// Close closes both branch sinks if they are non-nil.
func (c *ConditionalSink) Close() error {
	var first error
	if c.thenSink != nil {
		if err := c.thenSink.Close(); err != nil && first == nil {
			first = err
		}
	}
	if c.elseSink != nil {
		if err := c.elseSink.Close(); err != nil && first == nil {
			first = err
		}
	}
	return first
}
