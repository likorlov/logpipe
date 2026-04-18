/*
Package sink provides QuotaSink, which enforces a maximum number of log
entries per logical key within a sliding time window.

# Usage

	s := sink.NewQuotaSink(
		inner,
		100,              // max entries
		time.Minute,      // per window
		func(e logpipe.Entry) string {
			return e.Fields["service"] // quota per service
		},
	)

Entries that exceed the quota return an error and are not forwarded to the
wrapped sink. The counter resets automatically after the window elapses.
*/
package sink
