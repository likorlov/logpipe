package sink

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/example/logpipe"
)

// FileSink writes log entries as JSON lines to a file.
type FileSink struct {
	mu   sync.Mutex
	f    *os.File
	path string
}

// NewFileSink opens (or creates) the file at path for appending.
func NewFileSink(path string) (*FileSink, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("filesink: open %q: %w", path, err)
	}
	return &FileSink{f: f, path: path}, nil
}

// Write serialises the entry as a JSON line and appends it to the file.
func (s *FileSink) Write(entry logpipe.Entry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("filesink: marshal: %w", err)
	}
	data = append(data, '\n')

	s.mu.Lock()
	defer s.mu.Unlock()
	_, err = s.f.Write(data)
	if err != nil {
		return fmt.Errorf("filesink: write to %q: %w", s.path, err)
	}
	return nil
}

// Close flushes and closes the underlying file.
func (s *FileSink) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.f.Sync(); err != nil {
		return fmt.Errorf("filesink: sync %q: %w", s.path, err)
	}
	if err := s.f.Close(); err != nil {
		return fmt.Errorf("filesink: close %q: %w", s.path, err)
	}
	return nil
}

// Path returns the file path this sink writes to.
func (s *FileSink) Path() string { return s.path }

// Rotate closes the current file, renames it to path+".old", and opens a
// fresh file at the original path. It is safe to call concurrently with Write.
func (s *FileSink) Rotate() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.f.Sync(); err != nil {
		return fmt.Errorf("filesink: rotate sync %q: %w", s.path, err)
	}
	if err := s.f.Close(); err != nil {
		return fmt.Errorf("filesink: rotate close %q: %w", s.path, err)
	}
	if err := os.Rename(s.path, s.path+".old"); err != nil {
		return fmt.Errorf("filesink: rotate rename %q: %w", s.path, err)
	}
	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("filesink: rotate open %q: %w", s.path, err)
	}
	s.f = f
	return nil
}
