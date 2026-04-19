package sink_test

import (
	"os"

	"github.com/your/logpipe"
	"github.com/your/logpipe/sink"
)

func ExampleNewSequenceSink() {
	console := sink.NewConsoleSink(os.Stdout, false)
	s := sink.NewSequenceSink(console, "seq")
	defer s.Close()

	logger := logpipe.New(s)
	logger.Info("first message")
	logger.Info("second message")

	// Each entry will contain a "seq" field with values "1" and "2".
}
