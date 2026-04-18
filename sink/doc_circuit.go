/*
Package sink provides a CircuitBreaker sink that wraps any Sink and protects
downstream systems from repeated failures.

# Circuit Breaker

NewCircuitSink opens the circuit after a configurable number of consecutive
errors, rejecting writes with ErrCircuitOpen until the cooldown period elapses.
After cooldown, the next write attempt is forwarded to the inner sink; a
successful write resets the failure counter and closes the circuit again.

	inner := sink.NewWebhookSink("https://logs.example.com/ingest")
	s := sink.NewCircuitSink(inner, 5, 30*time.Second)
	defer s.Close()

	if err := s.Write(entry); err != nil {
		if errors.Is(err, sink.ErrCircuitOpen) {
			// circuit is open — downstream unavailable
		}
	}
*/
package sink
