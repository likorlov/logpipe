// Package sink provides output sink implementations for logpipe.
//
// Available sinks:
//
//   - ConsoleSink  – writes structured or pretty-printed log entries to an
//     io.Writer (e.g. os.Stdout).
//
//   - FileSink – appends JSON-encoded log entries to a file on disk.
//
//   - RotatingFileSink – like FileSink but rotates the output file once it
//     exceeds a configurable size limit.
//
//   - WebhookSink – POSTs each log entry as a JSON payload to an HTTP
//     endpoint, useful for forwarding logs to external services.
//
// All sinks implement the logpipe.Sink interface:
//
//	type Sink interface {
//	    Write(Entry) error
//	    Close() error
//	}
//
// Sinks can be composed via logpipe.New to fan-out entries to multiple
// destinations simultaneously.
package sink
