package sink

import (
	"sync"

	"github.com/logpipe/logpipe"
)

// ShadowSink forwards every entry to the primary sink and also sends a copy to
// the shadow sink. Errors from the shadow are silently discarded so they never
// affect the primary write path. This is useful for dark-launching a new sink
// alongside an existing one without risking production traffic.
type shadowSink struct {
	mu      sync.Mutex
	primary logpipe.Sink
	shadow  logpipe.Sink
}

// NewShadowSink returns a Sink that writes every entry to primary and, in a
// best-effort manner, also to shadow. Errors from shadow are ignored.
func NewShadowSink(primary, shadow logpipe.Sink) logpipe.Sink {
	if primary == nil {
		panic("logpipe/sink: NewShadowSink primary must not be nil")
	}
	if shadow == nil {
		panic("logpipe/sink: NewShadowSink shadow must not be nil")
	}
	return &shadowSink{primary: primary, shadow: shadow}
}

func (s *shadowSink) Write(entry logpipe.Entry) error {
	s.mu.Lock()
	primary := s.primary
	shadow := s.shadow
	s.mu.Unlock()

	// Shadow receives a shallow copy so mutations in primary do not affect it.
	copy := make(logpipe.Entry, len(entry))
	for k, v := range entry {
		copy[k] = v
	}
	// Best-effort: ignore shadow errors.
	_ = shadow.Write(copy)

	return primary.Write(entry)
}

func (s *shadowSink) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Close both; return primary error if any.
	_ = s.shadow.Close()
	return s.primary.Close()
}
