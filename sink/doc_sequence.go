// Package sink provides pluggable output sinks for logpipe.
//
// # SequenceSink
//
// SequenceSink wraps any Sink and injects a monotonically increasing sequence
// number into each log entry. This is useful for detecting dropped or
// reordered log messages downstream.
//
// Usage:
//
//	inner := sink.NewConsoleSink(os.Stdout, false)
//	s := sink.NewSequenceSink(inner, "seq")
//
//	logger := logpipe.New(s)
//	logger.Info("first")   // fields: {seq: "1"}
//	logger.Info("second")  // fields: {seq: "2"}
//
// The field name defaults to "seq" when an empty string is provided.
// Sequence numbers start at 1 and are safe for concurrent use.
package sink
