// Package sink provides pluggable output sinks for logpipe.
//
// # SamplingSink
//
// SamplingSink wraps any Sink and probabilistically forwards log entries
// based on a configured sample rate. This is useful for high-throughput
// services where only a fraction of debug or info logs need to be stored.
//
// Usage:
//
//	base := sink.NewConsoleSink(os.Stdout, false)
//	// Forward ~10 % of entries
//	sampled := sink.NewSamplingSink(base, 0.1, rand.NewSource(time.Now().UnixNano()))
//
//	logger := logpipe.New(logpipe.InfoLevel, sampled)
//	logger.Info("frequent event", nil)
//
// The sample rate is clamped to [0.0, 1.0]. A rate of 1.0 behaves like
// the underlying sink with no overhead beyond a single float64 comparison.
// A nil rand.Source falls back to a deterministic seed (42).
package sink
