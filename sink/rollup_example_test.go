package sink_test

import (
	"fmt"
	"time"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func ExampleNewRollupSink() {
	// Collect entries into a simple slice for the example.
	var received []logpipe.Entry
	collector := &funcSink{fn: func(e logpipe.Entry) error {
		received = append(received, e)
		return nil
	}}

	rollup := sink.NewRollupSink(collector, "duration_ms", 10*time.Second)

	// Simulate three request durations.
	_ = rollup.Write(logpipe.Entry{Fields: map[string]interface{}{"duration_ms": 120.0}})
	_ = rollup.Write(logpipe.Entry{Fields: map[string]interface{}{"duration_ms": 80.0}})
	_ = rollup.Write(logpipe.Entry{Fields: map[string]interface{}{"duration_ms": 200.0}})

	// Close flushes the current window immediately.
	_ = rollup.Close()

	if len(received) == 1 {
		e := received[0]
		fmt.Printf("count=%v sum=%v min=%v max=%v\n",
			e.Fields["count"], e.Fields["sum"],
			e.Fields["min"], e.Fields["max"])
	}
	// Output:
	// count=3 sum=400 min=80 max=200
}

// funcSink is a helper that wraps a plain function as a logpipe.Sink.
type funcSink struct {
	fn func(logpipe.Entry) error
}

func (f *funcSink) Write(e logpipe.Entry) error { return f.fn(e) }
func (f *funcSink) Close() error                { return nil }
