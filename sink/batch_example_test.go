package sink_test

import (
	"fmt"
	"time"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func ExampleNewBatchSink() {
	var received int
	flushFn := func(batch []logpipe.Entry) error {
		received += len(batch)
		return nil
	}

	s := sink.NewBatchSink(10, time.Second, flushFn)

	_ = s.Write(logpipe.Entry{Message: "a"})
	_ = s.Write(logpipe.Entry{Message: "b"})
	_ = s.Close() // flushes remaining entries

	fmt.Println(received)
	// Output: 2
}

func ExampleNewBatchSink_flushOnCapacity() {
	var flushCount int
	flushFn := func(batch []logpipe.Entry) error {
		flushCount++
		return nil
	}

	// batch size of 2 triggers flush after every 2 entries
	s := sink.NewBatchSink(2, time.Second, flushFn)

	_ = s.Write(logpipe.Entry{Message: "a"})
	_ = s.Write(logpipe.Entry{Message: "b"}) // triggers flush
	_ = s.Write(logpipe.Entry{Message: "c"})
	_ = s.Close() // flushes remaining entry

	fmt.Println(flushCount)
	// Output: 2
}
