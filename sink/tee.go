package sink

import (
	"github.com/logpipe/logpipe"
)

// TeeSink writes each log entry to a primary sink and, if the primary
// succeeds, also to a secondary sink — similar to the Unix tee command.
// Errors from the secondary sink are returned but do not prevent the
// primary write from being recorded.
type teeSink struct {
	primary   logpipe.Sink
	secondary logpipe.Sink
}

// NewTeeSink returns a Sink that writes every entry to primary and, on
// success, also to secondary. If primary fails the entry is not forwarded
// to secondary and the primary error is returned. If secondary fails its
// error is returned (primary write already happened).
func NewTeeSink(primary, secondary logpipe.Sink) logpipe.Sink {
	return &teeSink{primary: primary, secondary: secondary}
}

func (t *teeSink) Write(entry logpipe.Entry) error {
	if err := t.primary.Write(entry); err != nil {
		return err
	}
	return t.secondary.Write(entry)
}

func (t *teeSink) Close() error {
	pe := t.primary.Close()
	se := t.secondary.Close()
	if pe != nil {
		return pe
	}
	return se
}
