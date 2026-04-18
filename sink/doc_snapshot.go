/*
Package sink provides SnapshotSink, an in-memory sink that retains the last N
log entries for inspection, testing, and diagnostics.

# SnapshotSink

SnapshotSink stores up to a configurable maximum number of entries. When the
buffer is full the oldest entry is evicted (ring-buffer behaviour). It is
thread-safe and exposes Entries(), Len(), and Reset() helpers.

Example usage:

	snap := sink.NewSnapshotSink(50)
	logger := logpipe.New(snap)
	logger.Info("something happened")

	for _, e := range snap.Entries() {
		fmt.Println(e.Message)
	}
*/
package sink
