package sink_test

import (
	"fmt"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func ExampleNewScopeSink() {
	console := sink.NewConsoleSink(false)
	scoped := sink.NewScopeSink(console, "billing", "scope")

	_ = scoped.Write(logpipe.Entry{
		Level:   logpipe.INFO,
		Message: "invoice created",
		Fields:  map[string]any{"invoice_id": "INV-001"},
	})
	// Output entry will include Fields["scope"] = "billing"
	fmt.Println("scope injected")
	_ = scoped.Close()
	// Output:
	// scope injected
}
