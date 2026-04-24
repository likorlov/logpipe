/*
Package sink — CorrelateSink

CorrelateSink injects a correlation ID into every log entry it processes,
making it easy to trace chains of related events across a request, job,
or any other logical unit of work.

# Basic usage

	s := sink.NewCorrelateSink(inner, "", nil)
	// Writes will carry correlation_id="corr-1" automatically.
	_ = s.Write(logpipe.Entry{"msg": "starting job"})

# Custom ID generator

	import "github.com/google/uuid"

	s := sink.NewCorrelateSink(inner, "request_id", uuid.NewString)

# Rotating the ID

Call Rotate to start a new correlation scope (e.g. at the beginning of
each HTTP request):

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.Rotate()
		// all log entries within this handler share the new ID
	})

Entries that already carry the field are forwarded unchanged, so
downstream code can override the injected value when needed.
*/
package sink
