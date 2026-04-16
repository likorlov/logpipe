package sink_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/example/logpipe"
	"github.com/example/logpipe/sink"
)

func makeEntry(level logpipe.Level, msg string) logpipe.Entry {
	return logpipe.Entry{
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Level:     level,
		Message:   msg,
		Fields:    map[string]any{"key": "value"},
	}
}

func TestConsoleSink_Write(t *testing.T) {
	var buf bytes.Buffer
	s := sink.NewConsoleSink(&buf, false)

	entry := makeEntry(logpipe.INFO, "hello")
	if err := s.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := buf.String()
	if !strings.Contains(line, "hello") {
		t.Errorf("expected message in output, got: %s", line)
	}
	if !strings.HasSuffix(strings.TrimSpace(line), "}") {
		t.Errorf("expected JSON object, got: %s", line)
	}

	var out map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(line)), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
}

func TestConsoleSink_Pretty(t *testing.T) {
	var buf bytes.Buffer
	s := sink.NewConsoleSink(&buf, true)
	if err := s.Write(makeEntry(logpipe.WARN, "pretty")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "\n  ") {
		t.Error("expected indented JSON for pretty mode")
	}
}

func TestConsoleSink_Close(t *testing.T) {
	s := sink.NewConsoleSink(nil, false)
	if err := s.Close(); err != nil {
		t.Errorf("Close should be a no-op, got: %v", err)
	}
}
