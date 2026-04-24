/*
Package sink — FanoutSink

FanoutSink dispatches each log entry to all registered sinks concurrently,
waiting for every write to complete before returning. This is distinct from
MultiSink, which writes sequentially.

Use FanoutSink when the downstream sinks have non-trivial latency (e.g. network
calls) and you want to minimise the total wall-clock time spent per log entry.

All errors encountered across the concurrent writes are collected and returned
as a single combined error string so no failure is silently dropped.

	fs := sink.NewFanoutSink(
		sink.NewFileSink("/var/log/app.log"),
		sink.NewWebhookSink("https://logs.example.com/ingest"),
		sink.NewConsoleSink(),
	)
	defer fs.Close()

	logger := logpipe.New(fs)
	logger.Info("server started", logpipe.Fields{"port": 8080})
*/
package sink
