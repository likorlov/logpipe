package sink_test

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/logpipe"
	"github.com/example/logpipe/sink"
)

func TestRotatingFileSink_NoRotation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.log")

	s, err := sink.NewRotatingFileSink(path, 1024*1024)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	for i := 0; i < 3; i++ {
		if err := s.Write(logpipe.Entry{Level: logpipe.LevelInfo, Message: "msg", Time: time.Now()}); err != nil {
			t.Fatalf("Write %d: %v", i, err)
		}
	}
	s.Close()

	f, _ := os.Open(path)
	defer f.Close()
	sc := bufio.NewScanner(f)
	lines := 0
	for sc.Scan() {
		lines++
	}
	if lines != 3 {
		t.Errorf("want 3 lines, got %d", lines)
	}
}

func TestRotatingFileSink_Rotates(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.log")

	// tiny max to force rotation after first write
	s, err := sink.NewRotatingFileSink(path, 1)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	for i := 0; i < 3; i++ {
		if err := s.Write(logpipe.Entry{Level: logpipe.LevelWarn, Message: "rotate me", Time: time.Now()}); err != nil {
			t.Fatalf("Write %d: %v", i, err)
		}
		time.Sleep(time.Millisecond) // ensure unique timestamps
	}
	s.Close()

	entries, _ := filepath.Glob(filepath.Join(dir, "app.log*"))
	if len(entries) < 2 {
		t.Errorf("expected at least 2 files after rotation, got %d: %v", len(entries), entries)
	}
}
