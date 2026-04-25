package sink

import (
	"fmt"
	"sync"

	"github.com/logpipe/logpipe"
)

// CeilingSink forwards entries to inner only when the numeric value of a
// specified field is at or below a configurable ceiling. Entries whose field
// exceeds the ceiling, is missing, or is non-numeric are dropped.
//
// This is useful for suppressing runaway metrics or capping log verbosity
// based on a numeric severity, count, or rate field.
type ceilingSink struct {
	inner   logpipe.Sink
	field   string
	ceiling float64
	mu      sync.Mutex
}

// NewCeilingSink returns a Sink that forwards entries to inner only when the
// numeric value stored in field is less than or equal to ceiling. Entries that
// exceed the ceiling are silently dropped.
//
// Panics if ceiling is negative or inner is nil.
func NewCeilingSink(inner logpipe.Sink, field string, ceiling float64) logpipe.Sink {
	if inner == nil {
		panic("logpipe/sink: NewCeilingSink: inner sink must not be nil")
	}
	if ceiling < 0 {
		panic("logpipe/sink: NewCeilingSink: ceiling must be non-negative")
	}
	if field == "" {
		field = "value"
	}
	return &ceilingSink{inner: inner, field: field, ceiling: ceiling}
}

func (s *ceilingSink) Write(entry logpipe.Entry) error {
	s.mu.Lock()
	v, ok := entry[s.field]
	ceiling := s.ceiling
	s.mu.Unlock()

	if !ok {
		return nil
	}

	f, err := toFloat(v)
	if err != nil {
		return nil
	}

	if f > ceiling {
		return nil
	}

	return s.inner.Write(entry)
}

func (s *ceilingSink) Close() error {
	return s.inner.Close()
}

// toFloat converts common numeric types to float64.
func toFloat(v interface{}) (float64, error) {
	switch n := v.(type) {
	case float64:
		return n, nil
	case float32:
		return float64(n), nil
	case int:
		return float64(n), nil
	case int64:
		return float64(n), nil
	case int32:
		return float64(n), nil
	case uint:
		return float64(n), nil
	case uint64:
		return float64(n), nil
	}
	return 0, fmt.Errorf("not numeric")
}
