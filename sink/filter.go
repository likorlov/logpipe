package sink

import (
	"github.com/logpipe/logpipe"
)

// FilterFunc is a predicate that returns true if the entry should be passed
// through to the wrapped sink.
type FilterFunc func(entry logpipe.Entry) bool

// FilterSink wraps another Sink and only forwards entries that satisfy the
// provided predicate.
type FilterSink struct {
	inner  logpipe.Sink
	filter FilterFunc
}

// NewFilterSink creates a FilterSink that forwards entries to inner only when
// filter returns true.
func NewFilterSink(inner logpipe.Sink, filter FilterFunc) *FilterSink {
	if inner == nil {
		panic("logpipe/sink: NewFilterSink: inner sink must not be nil")
	}
	if filter == nil {
		panic("logpipe/sink: NewFilterSink: filter func must not be nil")
	}
	return &FilterSink{inner: inner, filter: filter}
}

// Write forwards entry to the inner sink only when the filter predicate
// returns true. Filtered-out entries are silently dropped.
func (f *FilterSink) Write(entry logpipe.Entry) error {
	if !f.filter(entry) {
		return nil
	}
	return f.inner.Write(entry)
}

// Close closes the inner sink.
func (f *FilterSink) Close() error {
	return f.inner.Close()
}

// LevelFilter returns a FilterFunc that passes only entries whose level is
// greater than or equal to min.
func LevelFilter(min logpipe.Level) FilterFunc {
	return func(entry logpipe.Entry) bool {
		return entry.Level >= min
	}
}

// FieldFilter returns a FilterFunc that passes only entries that contain the
// specified key in their Fields map.
func FieldFilter(key string) FilterFunc {
	return func(entry logpipe.Entry) bool {
		_, ok := entry.Fields[key]
		return ok
	}
}

// AndFilter returns a FilterFunc that passes only entries that satisfy all
// of the provided predicates. If no predicates are given, all entries pass.
func AndFilter(filters ...FilterFunc) FilterFunc {
	return func(entry logpipe.Entry) bool {
		for _, f := range filters {
			if !f(entry) {
				return false
			}
		}
		return true
	}
}
