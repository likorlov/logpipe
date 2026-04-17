/*
Package sink — RetrySink

RetrySink wraps any Sink and automatically retries failed Write calls
using a configurable number of attempts and an optional exponential backoff.

Basic usage:

	base := sink.NewWebhookSink("https://logs.example.com/ingest")
	s := sink.NewRetrySink(base, sink.RetryOptions{
		MaxAttempts: 5,
		Delay:       100 * time.Millisecond,
		Multiplier:  2.0, // 100ms, 200ms, 400ms, 800ms …
	})

Fields:

	- MaxAttempts  Total number of Write attempts (must be >= 1).
	- Delay        Initial wait between attempts. Zero means no sleep.
	- Multiplier   Factor applied to Delay after each failure.
	               1.0 = constant delay; 2.0 = exponential backoff.

If all attempts fail the last error is wrapped and returned to the caller.
Close is delegated directly to the wrapped Sink.
*/
package sink
