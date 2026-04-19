package sink_test

import (
	"testing"

	"github.com/your/logpipe"
	"github.com/your/logpipe/sink"
)

func BenchmarkSequenceSink_Write(b *testing.B) {
	col := &collectSink{}
	s := sink.NewSequenceSink(col, "seq")
	entry := logpipe.Entry{
		Message: "bench",
		Fields:  map[string]any{"key": "value"},
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = s.Write(entry)
		}
	})
}
