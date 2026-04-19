package sink_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func BenchmarkTimeoutSink_Write(b *testing.B) {
	inner := sink.NewConsoleSink(false)
	s := sink.NewTimeoutSink(inner, time.Second)
	defer s.Close()

	entry := logpipe.Entry{
		Level:   logpipe.Info,
		Message: "benchmark",
		Fields:  map[string]any{"k": "v"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Write(entry)
	}
}
