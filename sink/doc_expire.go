/*
Package sink provides ExpireSink, a sink that discards log entries whose
timestamp is older than a configured maximum age.

# Overview

ExpireSink reads a time.Time value from a named field in each log entry and
compares it against the current wall clock. If the entry is older than the
configured maxAge it is silently dropped; otherwise it is forwarded to the
inner sink unchanged.

Entries that do not carry the timestamp field, or whose field value is not a
time.Time, are always forwarded.

# Usage

	s := sink.NewExpireSink(
		inner,
		10*time.Minute,
		sink.WithExpireField("created_at"), // default: "ts"
	)

This is useful when replaying buffered or stored entries and you want to
automatically discard entries that are no longer relevant.
*/
package sink
