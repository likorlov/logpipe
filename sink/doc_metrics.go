/*
Package sink — MetricsSink

MetricsSink wraps any Sink and tracks three counters:

  - Writes: entries successfully forwarded to the underlying sink.
  - Drops:  entries rejected with ErrDropped (e.g. by a full async queue).
  - Errors: entries that caused any other non-nil error.

Usage:

	base := sink.NewConsoleSink(os.Stdout, false)
	ms  := sink.NewMetricsSink(base)

	logger := logpipe.New(logpipe.DEBUG, ms)
	logger.Info("hello", nil)

	fmt.Println(ms.Writes()) // 1

	// Reset counters between reporting windows.
	ms.Reset()
*/
package sink
