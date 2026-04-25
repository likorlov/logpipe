/*
Package sink — HedgeSink

HedgeSink reduces tail latency by issuing a speculative "hedge" write to a
secondary sink whenever the primary sink takes longer than a configured delay
to respond.

# Behaviour

  - If the primary sink responds within the delay, the secondary is never
    contacted and no extra goroutine overhead is incurred beyond the timer.
  - If the primary is still pending after the delay, the entry is concurrently
    written to the secondary sink. Whichever write succeeds first is treated as
    the result; the other is silently discarded.
  - If both writes fail, the first error returned is propagated to the caller.

# Example

	fast := sink.NewConsoleSink()           // local fallback
	slow := sink.NewWebhookSink(remoteURL)  // potentially slow remote

	h := sink.NewHedgeSink(slow, fast, 50*time.Millisecond)
	defer h.Close()

	logger := logpipe.New(h)
	logger.Info("request handled", "latency_ms", 12)

# Notes

The delay must be strictly positive; NewHedgeSink panics otherwise.
Both the primary and secondary sinks are closed when Close is called.
*/
package sink
