package sink_test

import (
	"fmt"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func ExampleNewReplaySink() {
	// Capture entries during a test, then replay into a real sink.
	rec := sink.NewReplaySink(nil)

	_ = rec.Write(logpipe.Entry{Level: logpipe.Info, Message: "startup"})
	_ = rec.Write(logpipe.Entry{Level: logpipe.Warn, Message: "low memory"})

	// Replay into a console sink for inspection.
	console := sink.NewConsoleSink(false)
	if err := rec.Replay(console); err != nil {
		fmt.Println("replay error:", err)
	}

	fmt.Println("recorded:", len(rec.Entries()))
	// Output:
	// recorded: 2
}
