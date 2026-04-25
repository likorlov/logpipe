package sink_test

import (
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func BenchmarkRegexSink_Match(b *testing.B) {
	col := &collectSink{}
	s, err := sink.NewRegexSink(col, "msg", `^error`, false)
	if err != nil {
		b.Fatal(err)
	}
	defer s.Close()
	e := logpipe.Entry{Fields: map[string]any{"msg": "error: something bad"}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		col.entries = col.entries[:0]
		_ = s.Write(e)
	}
}

func BenchmarkRegexSink_NoMatch(b *testing.B) {
	col := &collectSink{}
	s, _ := sink.NewRegexSink(col, "msg", `^error`, false)
	defer s.Close()
	e := logpipe.Entry{Fields: map[string]any{"msg": "info: all clear"}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Write(e)
	}
}
