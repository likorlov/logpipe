/*
Package sink provides OnceSink, a sink decorator that forwards only the
first log entry satisfying an optional predicate, dropping all subsequent
matches.

# Usage

	inner := sink.NewConsoleSink(os.Stderr, false)

	// Fire once on the first ERROR entry, ignore the rest.
	pred := func(e logpipe.Entry) bool { return e.Level == logpipe.ERROR }
	s := sink.NewOnceSink(inner, pred)

	// Reset so it can fire once more (e.g. after an alert is acknowledged).
	s.Reset()

OnceSink is safe for concurrent use.
*/
package sink
