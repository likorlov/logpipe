/*
Package sink provides NewTopKSink, a sink that tracks the top-K most
frequent string values for a specified entry field.

# Overview

NewTopKSink wraps any inner Sink and maintains a frequency table for a
chosen field. Every entry is forwarded to the inner sink unchanged. Call
TopK() at any time to retrieve the current leaderboard sorted by
descending count.

# Usage

	inner := sink.NewConsoleSink(os.Stdout, false)
	s := sink.NewTopKSink(inner, "status", 5)

	// ... write entries ...

	for _, e := range s.TopK() {
		fmt.Printf("%s: %d\n", e.Value(), e.Count())
	}

	// Reset counters between reporting windows.
	s.Reset()

# Notes

  - Only string field values are counted; other types are silently ignored.
  - TopK is safe for concurrent use.
  - k must be > 0 or NewTopKSink panics.
*/
package sink
