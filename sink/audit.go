package sink

import (
	"fmt"
	"sync"

	"github.com/logpipe/logpipe"
)

// AuditSink wraps an inner Sink and records every entry that passes through,
// tagging it with a monotonically increasing audit sequence number stored in
// the configured field (default "_audit_seq").
//
// Unlike SnapshotSink, AuditSink never evicts entries — it is intended for
// compliance/debug scenarios where a complete ordered record is required.
type AuditSink struct {
	inner logpipe.Sink
	field string

	mu      sync.Mutex
	seq     uint64
	entries []logpipe.Entry
}

// NewAuditSink creates an AuditSink wrapping inner. If field is empty the
// default field name "_audit_seq" is used.
func NewAuditSink(inner logpipe.Sink, field string) *AuditSink {
	if field == "" {
		field = "_audit_seq"
	}
	return &AuditSink{inner: inner, field: field}
}

func (a *AuditSink) Write(e logpipe.Entry) error {
	a.mu.Lock()
	a.seq++
	seq := a.seq

	// Build enriched copy without mutating the caller's entry.
	copy := logpipe.Entry{
		Level:   e.Level,
		Message: e.Message,
		Fields:  make(map[string]any, len(e.Fields)+1),
	}
	for k, v := range e.Fields {
		copy.Fields[k] = v
	}
	copy.Fields[a.field] = fmt.Sprintf("%d", seq)
	a.entries = append(a.entries, copy)
	a.mu.Unlock()

	return a.inner.Write(copy)
}

// Entries returns a snapshot of all audited entries in arrival order.
func (a *AuditSink) Entries() []logpipe.Entry {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]logpipe.Entry, len(a.entries))
	copy(out, a.entries)
	return out
}

// Reset clears the audit log and resets the sequence counter.
func (a *AuditSink) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entries = nil
	a.seq = 0
}

func (a *AuditSink) Close() error { return a.inner.Close() }
