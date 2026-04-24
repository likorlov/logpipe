/*
Package sink — PipelineSink

PipelineSink chains a sequence of sinks so that each entry flows through every
stage in order. This is useful when you want a series of transformations or
filters applied one after another before the entry reaches its final
destination.

Basic usage:

	p := sink.NewPipelineSink(
		sink.NewTimestampSink(console, "ts"),
		sink.NewRedactSink(console, []string{"password", "token"}),
		sink.NewConsoleSink(false),
	)

	_ = p.Write(logpipe.Entry{"msg": "user login", "password": "s3cr3t"})

Difference from MultiSink:

MultiSink broadcasts the same entry to all sinks in parallel. PipelineSink
passes the entry sequentially — each stage sees the output of the previous
one. Combine the two to fan-out a fully-processed entry to multiple backends.

Error handling:

If any stage returns an error, the pipeline stops immediately and the error is
returned to the caller. Stages after the failing one are not invoked.
*/
package sink
