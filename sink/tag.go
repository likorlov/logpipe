package sink

import (
	"github.com/logpipe/logpipe"
)

// TagSink wraps a sink and injects static key-value fields into every log entry
// before forwarding it downstream. Useful for annotating entries with service
// name, environment, region, etc.
type TagSink struct {
	next logpipe.Sink
	tags map[string]any
}

// NewTagSink returns a TagSink that merges tags into each entry's Fields before
// writing to next. Tags do not overwrite existing fields with the same key.
func NewTagSink(next logpipe.Sink, tags map[string]any) *TagSink {
	copy := make(map[string]any, len(tags))
	for k, v := range tags {
		copy[k] = v
	}
	return &TagSink{next: next, tags: copy}
}

func (s *TagSink) Write(entry logpipe.Entry) error {
	merged := make(map[string]any, len(s.tags)+len(entry.Fields))
	for k, v := range s.tags {
		merged[k] = v
	}
	// Entry fields take precedence over tags.
	for k, v := range entry.Fields {
		merged[k] = v
	}
	entry.Fields = merged
	return s.next.Write(entry)
}

func (s *TagSink) Close() error {
	return s.next.Close()
}
