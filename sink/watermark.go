package sink

import (
	"fmt"
	"sync"

	"github.com/logpipe/logpipe"
)

// WatermarkSink forwards entries to inner only when the named numeric field
// crosses a configured high-water or low-water threshold. Entries that do not
// carry the field, or whose value is not a number, are forwarded unchanged.
//
// direction:
//   "high" – forward when value >= threshold (default)
//   "low"  – forward when value <= threshold
type watermarkSink struct {
	mu        sync.Mutex
	inner     logpipe.Sink
	field     string
	threshold float64
	high      bool
}

// WatermarkOption configures a WatermarkSink.
type WatermarkOption func(*watermarkSink)

// WatermarkHigh configures the sink to forward entries whose field value is
// greater than or equal to threshold (this is the default direction).
func WatermarkHigh() WatermarkOption {
	return func(w *watermarkSink) { w.high = true }
}

// WatermarkLow configures the sink to forward entries whose field value is
// less than or equal to threshold.
func WatermarkLow() WatermarkOption {
	return func(w *watermarkSink) { w.high = false }
}

// NewWatermarkSink returns a Sink that forwards log entries to inner only when
// the numeric value stored in field crosses threshold in the configured
// direction. Entries missing the field or holding a non-numeric value pass
// through unconditionally.
func NewWatermarkSink(inner logpipe.Sink, field string, threshold float64, opts ...WatermarkOption) logpipe.Sink {
	w := &watermarkSink{
		inner:     inner,
		field:     field,
		threshold: threshold,
		high:      true,
	}
	for _, o := range opts {
		o(w)
	}
	return w
}

func (w *watermarkSink) Write(e logpipe.Entry) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	v, ok := e.Fields[w.field]
	if ok {
		val, err := toFloat64(v)
		if err == nil {
			if w.high && val < w.threshold {
				return nil
			}
			if !w.high && val > w.threshold {
				return nil
			}
		}
	}

	return w.inner.Write(e)
}

func (w *watermarkSink) Close() error {
	return w.inner.Close()
}

// toFloat64 converts common numeric types to float64.
func toFloat64(v interface{}) (float64, error) {
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
	return 0, fmt.Errorf("not a number: %T", v)
}
