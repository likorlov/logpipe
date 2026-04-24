package sink

import (
	"fmt"

	"github.com/logpipe/logpipe"
)

// flattenSink recursively flattens nested map[string]any fields in a log entry
// into dot-separated top-level keys before forwarding to the inner sink.
type flattenSink struct {
	inner  logpipe.Sink
	sep    string
	prefix string
}

// NewFlattenSink returns a Sink that flattens nested map fields using sep as
// the key separator (e.g. "."). Nested maps are expanded recursively; all
// other value types are left unchanged. The original entry is never mutated.
//
// Example: {"http": {"status": 200}} becomes {"http.status": 200}.
func NewFlattenSink(inner logpipe.Sink, sep string) logpipe.Sink {
	if sep == "" {
		sep = "."
	}
	return &flattenSink{inner: inner, sep: sep}
}

func (s *flattenSink) Write(entry logpipe.Entry) error {
	flat := make(logpipe.Fields, len(entry.Fields))
	flattenFields(flat, entry.Fields, "", s.sep)
	entry.Fields = flat
	return s.inner.Write(entry)
}

func (s *flattenSink) Close() error {
	return s.inner.Close()
}

// flattenFields recursively walks src, writing dot-joined keys into dst.
func flattenFields(dst logpipe.Fields, src logpipe.Fields, prefix, sep string) {
	for k, v := range src {
		key := k
		if prefix != "" {
			key = fmt.Sprintf("%s%s%s", prefix, sep, k)
		}
		switch child := v.(type) {
		case map[string]any:
			flattenFields(dst, logpipe.Fields(child), key, sep)
		case logpipe.Fields:
			flattenFields(dst, child, key, sep)
		default:
			dst[key] = v
		}
	}
}
