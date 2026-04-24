/*
Package sink provides ReorderSink, a buffering sink that collects a fixed
number of log entries and forwards them to an inner sink sorted by a numeric
field (ascending).

# Usage

	s := sink.NewReorderSink(
		sink.NewConsoleSink(os.Stdout, false),
		"sequence", // field name holding a numeric value
		64,         // buffer capacity; flush triggers when full
	)
	defer s.Close() // flushes any remaining buffered entries

# Behaviour

  - Entries are buffered until the buffer reaches the configured capacity,
    at which point the entire buffer is sorted and forwarded to the inner sink.
  - Calling Close flushes any remaining buffered entries before closing the
    inner sink.
  - The sort is stable: entries with equal field values retain their insertion
    order.
  - Non-numeric field values sort as 0.
*/
package sink
