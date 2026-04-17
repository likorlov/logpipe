package sink_test

import (
	"fmt"
	"time"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func ExampleNewTransformSink() {
	console := sink.NewConsoleSink(false)

	// Enrich every entry with a static service field and redact passwords.
	enriched := sink.NewTransformSink(console, func(e logpipe.Entry) (logpipe.Entry, bool) {
		if e.Fields == nil {
			e.Fields = map[string]interface{}{}
		}
		delete(e.Fields, "password")
		e.Fields["service"] = "example-svc"
		return e, true
	})

	_ = enriched.Write(logpipe.Entry{
		Level:   logpipe.INFO,
		Message: "user logged in",
		Time:    time.Now(),
		Fields:  map[string]interface{}{"user": "alice", "password": "s3cr3t"},
	})

	_ = enriched.Close()
	fmt.Println("done")
	// Output: done
}
