package sink

import (
	"github.com/logpipe/logpipe"
)

// TransformFunc is a function that transforms a log entry.
// It returns the modified entry and a bool indicating whether to forward it.
type TransformFunc func(entry logpipe.Entry) (logpipe.Entry, bool)

// transformSink wraps a sink and applies a TransformFunc to each entry.
type transformSink struct {
	next      logpipe.Sink
	transform TransformFunc
}

// NewTransformSink returns a Sink that applies fn to every entry before
// forwarding it to next. If fn returns false the entry is silently dropped.
//
//	redact := sink.NewTransformSink(next, func(e logpipe.Entry) (logpipe.Entry, bool) {
//		delete(e.Fields, "password")
//		return e, true
//	})
func NewTransformSink(next logpipe.Sink, fn TransformFunc) logpipe.Sink {
	return &transformSink{next: next, transform: fn}
}

func (t *transformSink) Write(entry logpipe.Entry) error {
	modified, ok := t.transform(entry)
	if !ok {
		return nil
	}
	return t.next.Write(modified)
}

func (t *transformSink) Close() error {
	return t.next.Close()
}
