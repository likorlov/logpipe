// Package sink provides pluggable output sinks for logpipe.
//
// # FallbackSink
//
// FallbackSink wraps two sinks — a primary and a fallback — and ensures log
// entries are delivered even when the primary sink is unavailable.
//
// On each Write, the primary sink is tried first. If it returns an error the
// same entry is forwarded to the fallback sink. Both sinks are closed when
// Close is called.
//
// Example:
//
//	primary := sink.NewWebhookSink("https://logs.example.com/ingest")
//	backup := sink.NewFileSink("/var/log/app-backup.log")
//
//	s := sink.NewFallbackSink(primary, backup)
//	defer s.Close()
//
//	logger := logpipe.New(s)
//	logger.Info("important event", nil)
package sink
