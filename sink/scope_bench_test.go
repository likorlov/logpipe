package sink_test

import (
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func BenchmarkScopeSink_Write(b *testing.B) {
	col := &collectSink{}
	s := sink.NewScopeSink(col, "benchscope", "scope")
	entry := logpipe.Entry{
		Level:   logpipe.INFO,
		Message: "benchmark",
		Fields:  map[string]any{"key": "value"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Write(entry)
	}
}
