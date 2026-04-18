package sink_test

import (
	"testing"
	"time"

	"github.com/your/logpipe"
	"github.com/your/logpipe/sink"
)

func BenchmarkQuotaSink_UnderLimit(b *testing.B) {
	dev := &captureSink{}
	q := sink.NewQuotaSink(dev, b.N+1, time.Hour, nil)
	e := logpipe.Entry{Message: "bench", Fields: map[string]any{"k": "v"}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Write(e) //nolint:errcheck
	}
}

func BenchmarkQuotaSink_OverLimit(b *testing.B) {
	dev := &captureSink{}
	q := sink.NewQuotaSink(dev, 1, time.Hour, nil)
	e := logpipe.Entry{Message: "bench"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Write(e) //nolint:errcheck
	}
}
