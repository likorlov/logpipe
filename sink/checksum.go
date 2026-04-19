package sink

import (
	"crypto/sha256"
	"fmt"

	"github.com/logpipe/logpipe"
)

// NewChecksumSink wraps inner and injects a SHA-256 checksum of the serialised
// fields into each log entry before forwarding it.
//
// The checksum is stored under field (default: "_checksum").
func NewChecksumSink(inner logpipe.Sink, field string) logpipe.Sink {
	if field == "" {
		field = "_checksum"
	}
	return &checksumSink{inner: inner, field: field}
}

type checksumSink struct {
	inner logpipe.Sink
	field string
}

func (s *checksumSink) Write(e logpipe.Entry) error {
	out := make(logpipe.Entry, len(e)+1)
	for k, v := range e {
		out[k] = v
	}
	out[s.field] = checksum(e)
	return s.inner.Write(out)
}

func (s *checksumSink) Close() error {
	return s.inner.Close()
}

// checksum returns a short hex digest of the entry's key-value pairs.
func checksum(e logpipe.Entry) string {
	h := sha256.New()
	for k, v := range e {
		fmt.Fprintf(h, "%s=%v;", k, v)
	}
	return fmt.Sprintf("%x", h.Sum(nil))[:16]
}
