/*
Package sink — TimeoutSink

TimeoutSink wraps any sink and enforces a maximum duration for each
Write call. If the underlying sink blocks longer than the configured
timeout the write is abandoned and an error is returned.

	s := sink.NewTimeoutSink(
		sink.NewConsoleSink(false),
		200*time.Millisecond,
	)
	defer s.Close()

	_ = s.Write(logpipe.Entry{Message: "hello"})

This is particularly useful when wrapping network sinks (webhook,
batch-over-HTTP, etc.) where latency spikes could otherwise stall
the calling goroutine indefinitely.
*/
package sink
