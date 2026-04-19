package sink_test

import (
	"fmt"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func ExampleNewMaskSink() {
	col := &collectSink{}
	s := sink.NewMaskSink(col,
		sink.MaskOption{
			Field:      "card",
			KeepPrefix: 4,
			KeepSuffix: 4,
			Mask:       "--------",
		},
	)
	_ = s.Write(logpipe.Entry{
		"level": "info",
		"card":  "1234567890123456",
	})
	fmt.Println(col.entries[0]["card"])
	// Output: 1234--------3456
}
