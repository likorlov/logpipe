package sink

import (
	"errors"
	"fmt"

	"github.com/andyh/logpipe"
)

// SchemaRule defines a required field and an optional type constraint.
type SchemaRule struct {
	Field    string
	TypeName string // e.g. "string", "int", "float64", "bool"; empty means any type
}

// schemaSink drops entries that do not satisfy all required schema rules.
type schemaSink struct {
	inner logpipe.Sink
	rules []SchemaRule
}

// NewSchemaSink returns a Sink that validates each log entry against the
// provided rules before forwarding it to inner. Entries missing a required
// field, or whose field value does not match the expected type, are silently
// dropped. If rules is empty every entry is forwarded unchanged.
//
//	schema := sink.NewSchemaSink(console, []sink.SchemaRule{
//		{Field: "level", TypeName: "string"},
//		{Field: "msg",   TypeName: "string"},
//		{Field: "ts"},
//	})
func NewSchemaSink(inner logpipe.Sink, rules []SchemaRule) logpipe.Sink {
	return &schemaSink{inner: inner, rules: rules}
}

func (s *schemaSink) Write(entry logpipe.Entry) error {
	for _, r := range s.rules {
		v, ok := entry.Fields[r.Field]
		if !ok {
			return nil // drop
		}
		if r.TypeName != "" && !matchesType(v, r.TypeName) {
			return nil // drop
		}
	}
	return s.inner.Write(entry)
}

func (s *schemaSink) Close() error { return s.inner.Close() }

// matchesType reports whether v's dynamic type matches the named Go type.
func matchesType(v any, name string) bool {
	switch name {
	case "string":
		_, ok := v.(string)
		return ok
	case "int":
		_, ok := v.(int)
		return ok
	case "float64":
		_, ok := v.(float64)
		return ok
	case "bool":
		_, ok := v.(bool)
		return ok
	default:
		return false
	}
}

// ErrSchemaMismatch is returned by ValidateEntry when a rule is violated.
var ErrSchemaMismatch = errors.New("schema mismatch")

// ValidateEntry checks entry against rules and returns a descriptive error on
// the first violation, or nil if all rules pass. Useful for explicit
// validation outside of the sink pipeline.
func ValidateEntry(entry logpipe.Entry, rules []SchemaRule) error {
	for _, r := range rules {
		v, ok := entry.Fields[r.Field]
		if !ok {
			return fmt.Errorf("%w: missing required field %q", ErrSchemaMismatch, r.Field)
		}
		if r.TypeName != "" && !matchesType(v, r.TypeName) {
			return fmt.Errorf("%w: field %q expected type %s", ErrSchemaMismatch, r.Field, r.TypeName)
		}
	}
	return nil
}
