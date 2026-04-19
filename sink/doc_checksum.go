// Package sink — ChecksumSink
//
// ChecksumSink computes a SHA-256 digest of every log entry's fields and
// injects the first 16 hex characters into the entry under a configurable
// field name (default: "_checksum") before forwarding to the inner sink.
//
// This is useful for tamper-detection, deduplication hints, or audit trails
// where a lightweight fingerprint of each record is required.
//
// Usage:
//
//	s := sink.NewChecksumSink(inner, "")         // uses "_checksum"
//	s := sink.NewChecksumSink(inner, "integrity") // custom field name
//
// The original entry is never mutated.
package sink
