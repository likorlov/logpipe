/*
Package sink provides RoundRobinSink, which distributes log entries evenly
across a set of downstream sinks using a round-robin strategy.

# Usage

	s1 := sink.NewConsoleSink(os.Stdout, false)
	s2 := sink.NewConsoleSink(os.Stderr, false)
	rr, err := sink.NewRoundRobinSink(s1, s2)
	if err != nil {
		log.Fatal(err)
	}
	defer rr.Close()

	// Entries alternate between s1 and s2.
	rr.Write(logpipe.Entry{"msg": "hello"})
	rr.Write(logpipe.Entry{"msg": "world"})

RoundRobinSink is safe for concurrent use. Close propagates to all
underlying sinks and returns the first error encountered.
*/
package sink
