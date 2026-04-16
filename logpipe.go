package logpipe

import (
	"sync"
	"time"
)

// Level represents the severity of a log entry.
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Entry is a structured log record.
type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     Level             `json:"level"`
	Message   string            `json:"message"`
	Fields    map[string]any    `json:"fields,omitempty"`
}

// Sink is the interface that output destinations must implement.
type Sink interface {
	Write(entry Entry) error
	Close() error
}

// Logger aggregates log entries and fans them out to registered sinks.
type Logger struct {
	mu    sync.RWMutex
	sinks []Sink
	min   Level
}

// New creates a Logger with the given minimum level.
func New(min Level) *Logger {
	return &Logger{min: min}
}

// AddSink registers a new output sink.
func (l *Logger) AddSink(s Sink) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.sinks = append(l.sinks, s)
}

// Log emits an entry to all sinks if level >= minimum.
func (l *Logger) Log(level Level, msg string, fields map[string]any) error {
	if level < l.min {
		return nil
	}
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Level:     level,
		Message:   msg,
		Fields:    fields,
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	for _, s := range l.sinks {
		if err := s.Write(entry); err != nil {
			return err
		}
	}
	return nil
}

// Close shuts down all sinks.
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, s := range l.sinks {
		if err := s.Close(); err != nil {
			return err
		}
	}
	return nil
}
