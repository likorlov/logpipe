package sink_test

import (
	"os"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func ExampleNewTeeSink() {
	// Write every entry to stdout (pretty) AND to a file.
	console := sink.NewConsoleSink(os.Stdout, true)
	file, err := sink.NewFileSink("/tmp/example-tee.log")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	tee := sink.NewTeeSink(file, console)
	logger := logpipe.New(tee)

	logger.Info("tee example", logpipe.F{"env": "prod"})
	// Output is written to both /tmp/example-tee.log and stdout.
}
