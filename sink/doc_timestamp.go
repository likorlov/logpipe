// Package sink provides pluggable output sinks for logpipe.
//
// # TimestampSink
//
// TimestampSink enriches every log entry with a UTC timestamp before
// forwarding it to an inner sink.
//
// Usage:
//
//	inner := sink.NewConsoleSink(os.Stdout, false)
//	s := sink.NewTimestampSink(inner, "ts")
//
//	// Each entry written through s will have a "ts" field containing
//	// the current time in RFC3339Nano format.
//	s.Write(logpipe.Entry{
//		Level:   logpipe.Info,
//		Message: "service started",
//		Fields:  map[string]any{},
//	})
//
// Pass an empty string as the field name to use the default field "ts".
package sink
