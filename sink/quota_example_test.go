package sink_test

import (
	"fmt"
	"time"

	"github.com/your/logpipe"
	"github.com/your/logpipe/sink"
)

func ExampleNewQuotaSink() {
	console := sink.NewConsoleSink(false)

	// Allow at most 5 log entries per service per minute.
	q := sink.NewQuotaSink(
		console,
		5,
		time.Minute,
		func(e logpipe.Entry) string {
			if svc, ok := e.Fields["service"].(string); ok {
				return svc
			}
			return "unknown"
		},
	)
	defer q.Close()

	err := q.Write(logpipe.Entry{
		Level:   logpipe.LevelInfo,
		Message: "request received",
		Fields:  map[string]any{"service": "api"},
	})
	fmt.Println(err)
	// Output: <nil>
}
