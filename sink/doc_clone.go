// Package sink provides pluggable output sinks for logpipe.
//
// # CloneSink
//
// CloneSink deep-copies every [logpipe.Entry] before forwarding it to the
// wrapped sink. Use it whenever a downstream sink (such as [TransformSink] or
// [EnrichSink]) modifies the entry's Fields map and you need to guarantee that
// the caller's original entry remains unchanged.
//
// Example:
//
//	s := sink.NewCloneSink(sink.NewConsoleSink(false))
//	logger := logpipe.New(logpipe.Info, s)
//	logger.Info("safe to mutate downstream", nil)
package sink
