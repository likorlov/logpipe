package sink_test

import (
	"fmt"
	"os"

	"github.com/andybar2/logpipe"
	"github.com/andybar2/logpipe/sink"
)

func ExampleNewMetricsSink() {
	base := sink.NewConsoleSink(os.Stdout, false)
	ms := sink.NewMetricsSink(base)

	logger := logpipe.New(logpipe.DEBUG, ms)
	logger.Info("hello", nil)
	logger.Info("world", nil)

	fmt.Printf("writes=%d drops=%d errors=%d\n",
		ms.Writes(), ms.Drops(), ms.Errors())
	// Output:
	// writes=2 drops=0 errors=0
}
