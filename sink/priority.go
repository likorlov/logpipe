package sink

import (
	"fmt"

	"github.com/logpipe/logpipe"
)

// PrioritySink routes each log entry to the first sink whose minimum level
// is satisfied by the entry's level. Sinks are evaluated in the order they
// are registered via Add.
type PrioritySink struct {
	routes []priorityRoute
}

type priorityRoute struct {
	minLevel logpipe.Level
	sink     logpipe.Sink
}

// NewPrioritySink returns an empty PrioritySink. Register routes with Add
// before use.
func NewPrioritySink() *PrioritySink {
	return &PrioritySink{}
}

// Add registers sink as a destination for entries whose level is >= minLevel.
// Routes are evaluated in insertion order; the first match wins.
func (p *PrioritySink) Add(minLevel logpipe.Level, s logpipe.Sink) {
	p.routes = append(p.routes, priorityRoute{minLevel: minLevel, sink: s})
}

// Write forwards the entry to the first matching sink.
func (p *PrioritySink) Write(entry logpipe.Entry) error {
	for _, r := range p.routes {
		if entry.Level >= r.minLevel {
			return r.sink.Write(entry)
		}
	}
	return nil
}

// Close closes all registered sinks, collecting errors.
func (p *PrioritySink) Close() error {
	var errs []error
	for _, r := range p.routes {
		if err := r.sink.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("priority sink close errors: %v", errs)
	}
	return nil
}
