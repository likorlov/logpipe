package sink

import (
	"fmt"
	"strings"

	"github.com/logpipe/logpipe"
)

// MaskSink replaces portions of string field values with a fixed mask pattern,
// useful for partially obscuring values like credit-card numbers or tokens
// while keeping a recognisable prefix or suffix.
type MaskSink struct {
	inner  logpipe.Sink
	fields map[string]maskRule
}

type maskRule struct {
	keepPrefix int
	keepSuffix int
	mask       string
}

// MaskOption configures a field masking rule.
type MaskOption struct {
	Field      string
	KeepPrefix int
	KeepSuffix int
	Mask       string
}

// NewMaskSink returns a Sink that partially masks the configured string fields
// before forwarding each entry to inner.
func NewMaskSink(inner logpipe.Sink, opts ...MaskOption) *MaskSink {
	fields := make(map[string]maskRule, len(opts))
	for _, o := range opts {
		mask := o.Mask
		if mask == "" {
			mask = "****"
		}
		fields[o.Field] = maskRule{keepPrefix: o.KeepPrefix, keepSuffix: o.KeepSuffix, mask: mask}
	}
	return &MaskSink{inner: inner, fields: fields}
}

func (s *MaskSink) Write(e logpipe.Entry) error {
	masked := make(logpipe.Entry, len(e))
	for k, v := range e {
		masked[k] = v
	}
	for field, rule := range s.fields {
		val, ok := masked[field]
		if !ok {
			continue
		}
		str := fmt.Sprintf("%v", val)
		masked[field] = applyMask(str, rule)
	}
	return s.inner.Write(masked)
}

func applyMask(s string, r maskRule) string {
	var b strings.Builder
	pfx := r.keepPrefix
	sfx := r.keepSuffix
	if pfx+sfx >= len(s) {
		return s
	}
	b.WriteString(s[:pfx])
	b.WriteString(r.mask)
	if sfx > 0 {
		b.WriteString(s[len(s)-sfx:])
	}
	return b.String()
}

func (s *MaskSink) Close() error { return s.inner.Close() }
