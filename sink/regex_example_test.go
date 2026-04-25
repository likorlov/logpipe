package sink_test

import (
	"fmt"
	"os"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func ExampleNewRegexSink() {
	console := sink.NewConsoleSink(os.Stdout, false)

	// Forward only entries whose "msg" field starts with "error".
	s, err := sink.NewRegexSink(console, "msg", `^error`, false)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	_ = s.Write(logpipe.Entry{Fields: map[string]any{"msg": "error: disk full"}})
	_ = s.Write(logpipe.Entry{Fields: map[string]any{"msg": "info: running fine"}})

	fmt.Println("done")
	// Output:
	// {"msg":"error: disk full"}
	// done
}
