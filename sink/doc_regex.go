// Package sink provides pluggable output sinks for logpipe.
//
// # RegexSink
//
// NewRegexSink wraps an inner Sink and filters log entries based on
// whether a named string field matches a regular expression.
//
// Entries whose field is absent or whose value is not a string are
// forwarded unconditionally, so the sink degrades gracefully when
// the field is optional.
//
// The invert parameter reverses the logic: when true, only entries
// that do NOT match the pattern are forwarded. This is useful for
// suppressing noisy patterns rather than selecting them.
//
// Example – forward only error-level messages:
//
//	s, err := sink.NewRegexSink(inner, "msg", `^error`, false)
//
// Example – suppress health-check noise:
//
//	s, err := sink.NewRegexSink(inner, "path", `/healthz`, true)
package sink
