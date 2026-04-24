/*
Package sink — StashSink

StashSink stores a copy of each log entry under a caller-supplied key so that
entries can be retrieved by key at any time. It is useful for correlating
request/response log pairs, caching the last-seen state for a session, or
building simple in-memory audit trails.

Basic usage:

	ss := sink.NewStashSink(innerSink, func(e logpipe.Entry) string {
		if id, ok := e["request_id"].(string); ok {
			return id
		}
		return ""
	})

	// Write an entry — it is both stashed and forwarded to innerSink.
	_ = ss.Write(logpipe.Entry{"request_id": "abc-123", "status": 200})

	// Retrieve it later.
	entry, ok := ss.Get("abc-123")

	// Remove it when no longer needed.
	ss.Delete("abc-123")

If keyFn returns an empty string the entry is forwarded but not stashed.
Writing a second entry under the same key overwrites the previous value.
Close clears the stash and closes the inner sink if it implements io.Closer.
*/
package sink
