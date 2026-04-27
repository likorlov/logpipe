/*
Package sink provides a ZipkinSink that forwards log entries to a
Zipkin-compatible distributed tracing endpoint using the Zipkin v2 HTTP/JSON
API.

# Usage

	s := sink.NewZipkinSink("http://localhost:9411/api/v2/spans")
	defer s.Close()

	logger := logpipe.New(logpipe.InfoLevel, s)
	logger.Info(logpipe.Fields{"message": "user.login", "user_id": "42"})

# Options

  - WithZipkinNameField — choose which entry field becomes the span name
    (default: "message").
  - WithZipkinHTTPClient — supply a custom *http.Client (e.g. with TLS or
    custom timeouts).

# Notes

All entry fields are attached as Zipkin span tags. The timestamp is set to the
moment Write is called. No TraceID or SpanID is generated; this sink is
intended for lightweight log-to-trace bridging rather than full distributed
tracing.
*/
package sink
