/*
Package sink — ReplaySink

ReplaySink records every log entry it receives and allows replaying them
into any other Sink at any point in time.

Usage:

	// Record-only mode (no inner sink).
	rec := sink.NewReplaySink(nil)
	logger.Write(logpipe.Entry{Message: "hello"})

	// Later, replay into a file sink.
	fs, _ := sink.NewFileSink("/var/log/app.log")
	if err := rec.Replay(fs); err != nil {
		log.Fatal(err)
	}

	// Or forward live AND record.
	ws := sink.NewWebhookSink("https://example.com/logs")
	rec2 := sink.NewReplaySink(ws)

	// Inspect what was sent.
	entries := rec2.Entries()
	fmt.Println(len(entries))
*/
package sink
