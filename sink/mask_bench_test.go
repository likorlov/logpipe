package sink_test

import (
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func BenchmarkMaskSink_Write(b *testing.B) {
	col := &collectSink{}
	s := sink.NewMaskSink(col,
		sink.MaskOption{Field: "card", KeepPrefix: 4, KeepSuffix: 4},
		sink.MaskOption{Field: "token", KeepPrefix: 3, KeepSuffix: 3},
	)
	e := logpipe.Entry{
		"level": "info",
		"card":  "1234567890123456",
		"token": "abcdefghijklmno",
		"msg":   "payment processed",
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		col.entries = col.entries[:0]
		_ = s.Write(e)
	}
}
