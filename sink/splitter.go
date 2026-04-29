package sink

import (
	"fmt"

	"github.com/logpipe/logpipe"
)

// SplitterFunc decides which sink key to route an entry to.
type SplitterFunc func(entry logpipe.Entry) string

type splitterSink struct {
	routes map[string]logpipe.Sink
	fn     SplitterFunc
	defaultKey string
}

// NewSplitterSink routes each log entry to one of several named sinks based on
// the value returned by fn. If the returned key has no matching sink, the entry
// is forwarded to the sink registered under defaultKey. If defaultKey is empty
// or not registered, the entry is silently dropped.
//
//	 splitter := sink.NewSplitterSink(
//	     func(e logpipe.Entry) string {
//	         if lvl, ok := e["level"].(string); ok {
//	             return lvl
//	         }
//	         return "default"
//	     },
//	     "default",
//	     map[string]logpipe.Sink{
//	         "error":   errorFileSink,
//	         "default": consoleSink,
//	     },
//	 )
func NewSplitterSink(fn SplitterFunc, defaultKey string, routes map[string]logpipe.Sink) logpipe.Sink {
	if fn == nil {
		panic("splitter: fn must not be nil")
	}
	if routes == nil {
		routes = make(map[string]logpipe.Sink)
	}
	return &splitterSink{fn: fn, defaultKey: defaultKey, routes: routes}
}

func (s *splitterSink) Write(entry logpipe.Entry) error {
	key := s.fn(entry)
	if sink, ok := s.routes[key]; ok {
		return sink.Write(entry)
	}
	if s.defaultKey != "" {
		if sink, ok := s.routes[s.defaultKey]; ok {
			return sink.Write(entry)
		}
	}
	return nil
}

func (s *splitterSink) Close() error {
	var errs []error
	seen := make(map[logpipe.Sink]struct{})
	for _, sk := range s.routes {
		if _, done := seen[sk]; done {
			continue
		}
		seen[sk] = struct{}{}
		if err := sk.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("splitter close errors: %v", errs)
	}
	return nil
}
