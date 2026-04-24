/*
Package sink provides JournalSink, a sink decorator that stamps each log entry
with a monotonically increasing journal index before forwarding it to an inner
sink.

# Usage

	inner := sink.NewConsoleSink(os.Stdout, false)
	j := sink.NewJournalSink(inner, "") // uses default field "_journal"

	_ = j.Write(logpipe.Entry{"msg": "started"})
	_ = j.Write(logpipe.Entry{"msg": "processing"})

	fmt.Println(j.Index()) // 2

The journal index is injected into every outgoing entry under the configured
field name (default: "_journal"). The original entry is never mutated.

The index can be reset at any time via Reset(), which is useful when replaying
or restarting a processing pipeline.
*/
package sink
