/*
Package sink provides RateLimitSink, a sink decorator that enforces a
maximum number of log entries forwarded to a wrapped sink within a
configurable time interval.

# Usage

	base := sink.NewConsoleSink(os.Stdout, false)
	// allow at most 100 log entries per second
	rl := sink.NewRateLimitSink(base, 100, time.Second)

	logger := logpipe.New(logpipe.InfoLevel, rl)
	logger.Info("hello", nil)

# Behaviour

Entries that arrive after the per-interval quota has been exhausted are
silently dropped (Write returns nil). The quota resets automatically once
the interval elapses, starting from the moment the first entry after the
previous reset is received.
*/
package sink
