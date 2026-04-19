/*
Package sink provides ConditionalSink, which routes each log entry to one of
two downstream sinks based on a caller-supplied predicate.

# Usage

	errorSink := sink.NewConsoleSink(os.Stderr, false)
	infoSink  := sink.NewConsoleSink(os.Stdout, false)

	cs := sink.NewConditionalSink(
		func(e logpipe.Entry) bool { return e.Level >= logpipe.LevelError },
		errorSink,
		infoSink,
	)

	logger := logpipe.New(cs)
	logger.Info("hello")  // → infoSink
	logger.Error("oops") // → errorSink

Either branch sink may be nil; entries routed to a nil sink are silently
dropped without error. Close propagates to both non-nil sinks and returns the
first error encountered.
*/
package sink
