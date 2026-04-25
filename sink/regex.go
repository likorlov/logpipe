package sink

import (
	"fmt"
	"regexp"

	"github.com/logpipe/logpipe"
)

// regexSink forwards entries to inner only when a string field matches
// (or does not match) a compiled regular expression.
type regexSink struct {
	inner   logpipe.Sink
	field   string
	re      *regexp.Regexp
	invert  bool
}

// NewRegexSink returns a Sink that filters entries based on whether the
// value of field matches pattern. When invert is true the entry is
// forwarded only when the pattern does NOT match. Entries whose field is
// absent or not a string are always forwarded.
func NewRegexSink(inner logpipe.Sink, field, pattern string, invert bool) (logpipe.Sink, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("sink/regex: invalid pattern %q: %w", pattern, err)
	}
	return &regexSink{inner: inner, field: field, re: re, invert: invert}, nil
}

func (s *regexSink) Write(e logpipe.Entry) error {
	val, ok := e.Fields[s.field]
	if !ok {
		return s.inner.Write(e)
	}
	str, ok := val.(string)
	if !ok {
		return s.inner.Write(e)
	}
	matched := s.re.MatchString(str)
	if s.invert {
		matched = !matched
	}
	if !matched {
		return nil
	}
	return s.inner.Write(e)
}

func (s *regexSink) Close() error {
	return s.inner.Close()
}
