// Package sink provides AggregateSink, which buffers log entries and merges
// them into a single combined entry once a configurable batch size is reached.
//
// # Usage
//
//	inner := sink.NewConsoleSink(os.Stdout, false)
//	agg := sink.NewAggregateSink(inner, 50, "batch_size")
//	defer agg.Close()
//
//	// Write individual entries; every 50 entries are merged and forwarded.
//	agg.Write(entry)
//
// Fields from all buffered entries are merged into the outgoing entry; on key
// collision the latest entry wins. The number of merged entries is stored
// under the configured field name (default "agg_count").
//
// Call Close to flush any remaining buffered entries before shutdown.
package sink
