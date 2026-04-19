// Package sink provides pluggable output sinks for logpipe.
//
// # HashRoute Sink
//
// NewHashRouteSink shards log entries across multiple downstream sinks by
// computing a stable FNV-32a hash of a chosen field value.
//
// This is useful for partitioning log streams by tenant, user, or any other
// categorical field while guaranteeing that all entries sharing the same
// field value are always routed to the same downstream sink.
//
// Entries that do not contain the target field are always sent to the first
// sink (index 0).
//
// Example:
//
//	fileA, _ := sink.NewFileSink("/var/log/shard-0.log")
//	fileB, _ := sink.NewFileSink("/var/log/shard-1.log")
//
//	router := sink.NewHashRouteSink("tenant_id", fileA, fileB)
//	logger := logpipe.New(router)
//	logger.Info("request handled", map[string]any{"tenant_id": "acme"})
package sink
