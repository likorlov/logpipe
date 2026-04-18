package sink

import (
	"strings"

	"github.com/logpipe/logpipe"
)

// RedactSink masks sensitive field values before forwarding entries to the
// wrapped sink. Each matching field value is replaced with the given mask
// string (defaults to "***").
type redactSink struct {
	next   logpipe.Sink
	fields map[string]struct{}
	mask   string
}

// NewRedactSink returns a Sink that redacts the given field keys in every log
// entry before passing it downstream. If mask is empty, "***" is used.
func NewRedactSink(next logpipe.Sink, mask string, fields ...string) logpipe.Sink {
	if mask == "" {
		mask = "***"
	}
	fm := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		fm[strings.ToLower(f)] = struct{}{}
	}
	return &redactSink{next: next, fields: fm, mask: mask}
}

func (r *redactSink) Write(e logpipe.Entry) error {
	if len(r.fields) == 0 || len(e.Fields) == 0 {
		return r.next.Write(e)
	}
	// Shallow-copy fields map so we don't mutate the original entry.
	redacted := make(map[string]any, len(e.Fields))
	for k, v := range e.Fields {
		if _, ok := r.fields[strings.ToLower(k)]; ok {
			redacted[k] = r.mask
		} else {
			redacted[k] = v
		}
	}
	e.Fields = redacted
	return r.next.Write(e)
}

func (r *redactSink) Close() error {
	return r.next.Close()
}
