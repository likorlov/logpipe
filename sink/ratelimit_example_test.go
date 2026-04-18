package sink_test

import (
	"os"
	"time"

	"github.com/your/logpipe"
	"github.com/your/logpipe/sink"
)

func ExampleNewRateLimitSink() {
	base := sink.NewConsoleSink(os.Stdout, false)
	// Allow at most 2 log entries per second.
	rl := sink.NewRateLimitSink(base, 2, time.Second)
	defer rl.Close()

	logger := logpipe.New(logpipe.DebugLevel, rl)
	logger.Info("first", nil)
	logger.Info("second", nil)
	// third entry is dropped — quota exhausted
	logger.Info("third", nil)
}
