// Package sink provides pluggable output sinks for logpipe.
//
// # WindowSink
//
// WindowSink wraps an inner sink and maintains an in-memory sliding window of
// log entries. Every entry is forwarded to the inner sink immediately; the
// window is used for inspection via Entries().
//
// Entries older than the configured window duration are evicted automatically
// on each Write call.
//
// Example:
//
//	col := sink.NewConsoleSink()
//	w := sink.NewWindowSink(col, 30*time.Second)
//
//	logger.Write(logpipe.Entry{"msg": "hello"})
//
//	// inspect entries logged in the last 30 seconds
//	recent := w.Entries()
package sink
