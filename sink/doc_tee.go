/*
Package sink — TeeSink

TeeSink mirrors every log entry to two sinks: a primary and a secondary.
This is useful when you want a canonical destination (e.g. a file) and an
additional observer (e.g. a webhook or an in-memory buffer for testing)
without building a full fan-out pipeline.

Behaviour:
  - If the primary write fails the entry is NOT forwarded to the secondary
    and the primary error is returned.
  - If the primary write succeeds but the secondary write fails the
    secondary error is returned (the primary record is already committed).
  - Close flushes both sinks; the first non-nil error is returned.

Example:

	primary := sink.NewFileSink("/var/log/app.log")
	debug  := sink.NewConsoleSink(os.Stderr, false)

	tee := sink.NewTeeSink(primary, debug)
	logger := logpipe.New(tee)
	logger.Info("request handled", logpipe.F{"status": 200})
*/
package sink
