// Package sink provides pluggable output sinks for logpipe.
//
// # BufferedSink
//
// BufferedSink wraps any Sink and accumulates log entries in memory,
// flushing them to the inner sink either when the buffer reaches a
// configured size threshold or when a periodic interval elapses.
// This is useful for reducing I/O pressure on high-throughput sinks
// such as file or webhook sinks.
//
// Example:
//
//	inner := sink.NewFileSink("/var/log/app.log")
//	buffered := sink.NewBufferedSink(inner, 50, 5*time.Second)
//	defer buffered.Close()
//
// # MultiSink
//
// MultiSink fans log entries out to multiple sinks simultaneously.
// All sinks receive every entry; errors from individual sinks are
// collected and joined into a single combined error.
//
// Example:
//
//	console := sink.NewConsoleSink(true)
//	file, _ := sink.NewFileSink("/var/log/app.log")
//	multi := sink.NewMultiSink(console, file)
//	defer multi.Close()
package sink
