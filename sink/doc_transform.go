/*
Package sink provides a TransformSink that mutates or drops log entries
before forwarding them to a downstream sink.

# TransformSink

NewTransformSink wraps any Sink and applies a user-supplied TransformFunc to
every entry. The function receives a copy of the entry and returns the
(possibly modified) entry plus a boolean. When the boolean is false the entry
is silently discarded and the downstream sink is never called.

Common use-cases:
  - Redacting sensitive fields (passwords, tokens)
  - Enriching entries with additional context (hostname, version)
  - Normalising field names across services

Example:

	redact := sink.NewTransformSink(next, func(e logpipe.Entry) (logpipe.Entry, bool) {
		delete(e.Fields, "secret")
		e.Fields["service"] = "my-svc"
		return e, true
	})
*/
package sink
