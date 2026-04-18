package sink_test

import (
	"testing"
	"time"

	"github.com/your/logpipe"
	"github.com/your/logpipe/sink"
)

func BenchmarkSnapshotSink_Write(b *testing.B) {
	s := sink.NewSnapshotSink(1000)
	e := logpipe.Entry{
		Level:   logpipe.DEBUG,
		Message: "benchmark entry",
		Time:    time.Now(),
		Fields:  map[string]any{"key": "value"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Write(e)
	}
}

func BenchmarkSnapshotSink_Entries(b *testing.B) {
	s := sink.NewSnapshotSink(1000)
	e := logpipe.Entry{Level: logpipe.INFO, Message: "x", Time: time.Now()}
	for i := 0; i < 1000; i++ {
		s.Write(e)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Entries()
	}
}
