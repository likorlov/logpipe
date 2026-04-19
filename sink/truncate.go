package sink

import (
	"fmt"

	"github.com/yourusername/logpipe"
)

// TruncateSink truncates string field values that exceed a maximum length.
// Fields not present or not strings are passed through unchanged.
type TruncateSink struct {
	inner  logpipe.Sink
	maxLen int
	fields []string
	suffix string
}

// NewTruncateSink returns a sink that truncates the given fields to maxLen runes.
// If suffix is non-empty (e.g. "...") it is appended after truncation.
func NewTruncateSink(inner logpipe.Sink, maxLen int, suffix string, fields ...string) *TruncateSink {
	if maxLen <= 0 {
		panic("logpipe/sink: TruncateSink maxLen must be > 0")
	}
	return &TruncateSink{
		inner:  inner,
		maxLen: maxLen,
		fields: fields,
		suffix: suffix,
	}
}

func (s *TruncateSink) Write(entry logpipe.Entry) error {
	copy := make(logpipe.Entry, len(entry))
	for k, v := range entry {
		copy[k] = v
	}
	for _, f := range s.fields {
		v, ok := copy[f]
		if !ok {
			continue
		}
		str, ok := v.(string)
		if !ok {
			continue
		}
		runes := []rune(str)
		if len(runes) > s.maxLen {
			truncated := string(runes[:s.maxLen])
			copy[f] = fmt.Sprintf("%s%s", truncated, s.suffix)
		}
	}
	return s.inner.Write(copy)
}

func (s *TruncateSink) Close() error {
	return s.inner.Close()
}
