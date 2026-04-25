/*
Package sink — MergeSink

MergeSink deep-merges a fixed set of fields into every log entry before
forwarding it to the inner sink. It is similar to EnrichSink but handles
nested maps recursively: when both the base fields and the entry contain a
map[string]any under the same key, the two maps are merged rather than one
replacing the other.

Entry fields always win over the merge fields at every level of nesting.

Example usage:

	base := sink.NewConsoleSink(os.Stdout, false)

	// Attach deployment metadata to every log entry.
	s := sink.NewMergeSink(base, map[string]any{
		"env": "production",
		"meta": map[string]any{
			"service": "payments",
			"region":  "eu-west-1",
		},
	})

	// An entry that also carries a "meta" map will be deep-merged:
	//   meta.service comes from the entry, meta.region from the base.
	_ = s.Write(logpipe.Entry{
		Level:   logpipe.Info,
		Message: "charge processed",
		Fields:  map[string]any{"meta": map[string]any{"service": "checkout"}},
	})
*/
package sink
