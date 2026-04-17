/*
Package sink provides DedupeSink, a sink wrapper that suppresses repeated log
entries within a sliding time window.

# Overview

In high-throughput services the same error can be logged thousands of times per
second. DedupeSink transparently deduplicates entries that share the same level
and message, forwarding only the first occurrence within the configured window.

# Usage

	base := sink.NewConsoleSink(os.Stdout, false)
	deduped := sink.NewDedupeSink(base, 5*time.Second)

	logger := logpipe.New(logpipe.Info, deduped)
	logger.Error("database unreachable", nil) // forwarded
	logger.Error("database unreachable", nil) // suppressed for 5 s

# Notes

  - The deduplication key is (level, message); structured fields are ignored.
  - DedupeSink is safe for concurrent use.
  - The seen-entry map grows unboundedly; for very high cardinality messages
    consider pairing DedupeSink with a short window.
*/
package sink
