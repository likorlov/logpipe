package sink_test

import (
	"fmt"
	"time"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func ExampleNewTimeoutSink() {
	// Wrap a console sink so that each write must complete within 500ms.
	console := sink.NewConsoleSink(false)
	s := sink.NewTimeoutSink(console, 500*time.Millisecond)
	defer s.Close()

	err := s.Write(logpipe.Entry{
		Level:   logpipe.Info,
		Message: "processing request",
		Fields:  map[string]any{"request_id": "abc123"},
	})
	if err != nil {
		fmt.Println("write failed:", err)
	}
}
