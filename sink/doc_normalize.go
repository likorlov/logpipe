// Package sink provides pluggable output sinks for logpipe.
//
// # NormalizeSink
//
// NormalizeSink rewrites the field keys of every log entry using a
// caller-supplied normalization function before forwarding the entry to an
// inner sink.
//
// This is useful when log entries are produced by multiple subsystems with
// inconsistent key casing conventions and you want to enforce a uniform style
// (e.g. all-lowercase or snake_case) before the entries reach storage.
//
// Usage:
//
//	// Lowercase all field keys (default when normFn is nil)
//	s := sink.NewNormalizeSink(inner, nil)
//
//	// Custom normalizer — convert keys to UPPER_CASE
//	s := sink.NewNormalizeSink(inner, strings.ToUpper)
//
// The original entry is never mutated; a shallow copy with new keys is
// created for each write.
package sink
