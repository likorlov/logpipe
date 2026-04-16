package sink_test

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/example/logpipe"
	"github.com/example/logpipe/sink"
)

func tempFile(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp("", "logpipe-*.log")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestFileSink_Write(t *testing.T) {
	path := tempFile(t)
	s, err := sink.NewFileSink(path)
	if err != nil {
		t.Fatalf("NewFileSink: %v", err)
	}

	entry := logpipe.Entry{
		Level:   logpipe.LevelInfo,
		Message: "hello file",
		Time:    time.Now(),
		Fields:  map[string]any{"key": "value"},
	}
	if err := s.Write(entry); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	f, _ := os.Open(path)
	defer f.Close()
	var got logpipe.Entry
	if err := json.NewDecoder(bufio.NewReader(f)).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Message != entry.Message {
		t.Errorf("message: got %q want %q", got.Message, entry.Message)
	}
}

func TestFileSink_MultipleEntries(t *testing.T) {
	path := tempFile(t)
	s, _ := sink.NewFileSink(path)

	for i := 0; i < 5; i++ {
		_ = s.Write(logpipe.Entry{Level: logpipe.LevelDebug, Message: "line", Time: time.Now()})
	}
	s.Close()

	f, _ := os.Open(path)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		count++
	}
	if count != 5 {
		t.Errorf("expected 5 lines, got %d", count)
	}
}

func TestFileSink_BadPath(t *testing.T) {
	_, err := sink.NewFileSink("/no/such/dir/logpipe.log")
	if err == nil {
		t.Error("expected error for bad path")
	}
}
