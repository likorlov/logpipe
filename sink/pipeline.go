package sink

import "github.com/logpipe/logpipe"

// PipelineSink chains multiple sinks in sequence, passing the (potentially
// mutated) entry from one sink to the next. Unlike MultiSink, each sink in
// the pipeline receives the entry as modified by the previous sink.
//
// If any sink in the chain returns an error the pipeline stops and the error
// is returned immediately; subsequent sinks are not called.
//
// Close closes all sinks in order, collecting the first error encountered.
type PipelineSink struct {
	stages []logpipe.Sink
}

// NewPipelineSink creates a PipelineSink that passes each log entry through
// the provided stages in order. Panics if no stages are provided.
func NewPipelineSink(stages ...logpipe.Sink) *PipelineSink {
	if len(stages) == 0 {
		panic("logpipe/sink: NewPipelineSink requires at least one stage")
	}
	return &PipelineSink{stages: stages}
}

// Write passes the entry through each stage in order. If a stage returns a
// non-nil error the pipeline halts and that error is returned.
func (p *PipelineSink) Write(entry logpipe.Entry) error {
	current := entry
	for _, s := range p.stages {
		if err := s.Write(current); err != nil {
			return err
		}
	}
	return nil
}

// Close closes all stages in declaration order, returning the first error
// encountered. All stages are closed regardless of intermediate errors.
func (p *PipelineSink) Close() error {
	var first error
	for _, s := range p.stages {
		if err := s.Close(); err != nil && first == nil {
			first = err
		}
	}
	return first
}
