/*
Package sink — ShadowSink

ShadowSink enables dark-launching a new sink alongside an existing production
sink with zero risk. Every log entry is forwarded to the primary sink and a
best-effort copy is sent to the shadow sink. Errors from the shadow are
silently discarded so they can never affect the primary write path.

Usage:

	// Existing production sink.
	prod := sink.NewFileSink("/var/log/app.log")

	// New experimental sink under evaluation.
	experimental, _ := sink.NewWebhookSink("https://new-backend.example.com/logs")

	// Shadow traffic to the experimental sink without risk.
	s := sink.NewShadowSink(prod, experimental)

	logger := logpipe.New(s)
	logger.Info("request handled", "status", 200)

Both sinks receive every entry. If the experimental sink returns an error or
becomes unavailable, the logger continues operating normally through the
primary sink.
*/
package sink
