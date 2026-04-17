package sink_test

import (
	"fmt"
	"time"

	"github.com/your/logpipe"
	"github.com/your/logpipe/sink"
)

func ExampleNewRetrySink() {
	// Wrap a console sink (always succeeds) just to illustrate the API.
	base := sink.NewConsoleSink(false)
	s := sink.NewRetrySink(base, sink.RetryOptions{
		MaxAttempts: 3,
		Delay:       10 * time.Millisecond,
		Multiplier:  2.0,
	})
	defer s.Close()

	err := s.Write(logpipe.Entry{
		Level:   logpipe.Info,
		Message: "hello from retry sink",
		Fields:  map[string]any{"service": "demo"},
	})
	if err != nil {
		fmt.Println("error:", err)
	}
}
