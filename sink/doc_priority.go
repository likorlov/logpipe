/*
Package sink provides PrioritySink, which routes each log entry to the
first registered sink whose minimum level threshold is met.

# Usage

	p := sink.NewPrioritySink()
	p.Add(logpipe.LevelError, alertSink)  // errors go to alerting
	p.Add(logpipe.LevelDebug, consoleSink) // everything else to console

	logger := logpipe.New(p)

Routes are evaluated in insertion order. The first sink whose minLevel is
<= the entry level receives the entry; remaining sinks are skipped.
Entries that match no route are silently dropped.
*/
package sink
