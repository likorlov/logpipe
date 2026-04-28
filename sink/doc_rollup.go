/*
Package sink — RollupSink

RollupSink aggregates numeric field values over a fixed time window and emits
a single summary entry per window to the inner sink. It is useful for reducing
high-cardinality metric streams to periodic roll-ups without losing statistical
information.

Each flushed entry contains:

  - "field"  – the name of the monitored field
  - "count"  – number of entries received in the window (int64)
  - "sum"    – total of all values (float64)
  - "min"    – smallest value seen (float64)
  - "max"    – largest value seen (float64)

Entries whose monitored field is absent or non-numeric are silently dropped.
The final window is always flushed when Close is called.

Example:

	rollup := sink.NewRollupSink(console, "response_time_ms", time.Minute)
	defer rollup.Close()
	// Every minute one summary entry is forwarded to console.
*/
package sink
