/*
Package sink — CoalesceSink

CoalesceSink tries a list of sinks in order and returns on the first
successful write. Unlike FallbackSink (which is limited to a primary and
a single fallback), CoalesceSink accepts an arbitrary number of sinks,
making it suitable for multi-tier delivery chains.

Example:

	primary := sink.NewWebhookSink("https://logs.example.com/ingest")
	secondary := sink.NewFileSink("/var/log/app/fallback.log")
	tertiary := sink.NewConsoleSink(os.Stderr, false)

	s := sink.NewCoalesceSink(primary, secondary, tertiary)

	logger := logpipe.New(s)
	logger.Info("order fulfilled", map[string]any{"order_id": 42})

If primary succeeds the entry is not forwarded to secondary or tertiary.
If primary fails but secondary succeeds, tertiary is never tried. Only
when every sink returns an error is the combined error surfaced to the
caller.
*/
package sink
