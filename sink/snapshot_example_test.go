package sink_test

import (
	"fmt"
	"time"

	"github.com/your/logpipe"
	"github.com/your/logpipe/sink"
)

func ExampleNewSnapshotSink() {
	snap := sink.NewSnapshotSink(100)

	_ = snap.Write(logpipe.Entry{
		Level:   logpipe.INFO,
		Message: "service started",
		Time:    time.Now(),
	})
	_ = snap.Write(logpipe.Entry{
		Level:   logpipe.WARN,
		Message: "high memory usage",
		Time:    time.Now(),
	})

	fmt.Println(snap.Len())
	for _, e := range snap.Entries() {
		fmt.Println(e.Message)
	}
	// Output:
	// 2
	// service started
	// high memory usage
}
