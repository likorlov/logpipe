// Package sink provides output sink implementations for logpipe.
//
// Available sinks:
//
//   - ConsoleSink  – writes human-readable or JSON output to stdout/stderr.
//   - FileSink     – appends JSON-line entries to a single file.
//   - RotatingFileSink – like FileSink but rotates the file once it exceeds
//     a configurable size limit, renaming the old file with a timestamp suffix.
//
// All sinks implement the logpipe.Sink interface:
//
//	type Sink interface {
//	    Write(Entry) error
//	    Close() error
//	}
//
// Example – write to a rotating file:
//
//	s, err := sink.NewRotatingFileSink("/var/log/myapp.log", 50*1024*1024)
//	if err != nil { log.Fatal(err) }
//	logger := logpipe.New(logpipe.LevelInfo, s)
//	defer logger.Close()
package sink
