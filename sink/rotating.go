package sink

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/example/logpipe"
)

// RotatingFileSink writes JSON-line logs and rotates when the file exceeds MaxBytes.
type RotatingFileSink struct {
	mu       sync.Mutex
	base     string
	MaxBytes int64
	current  *FileSink
	size     int64
}

// NewRotatingFileSink creates a sink that rotates at maxBytes.
func NewRotatingFileSink(basePath string, maxBytes int64) (*RotatingFileSink, error) {
	if maxBytes <= 0 {
		maxBytes = 10 * 1024 * 1024 // 10 MiB default
	}
	s := &RotatingFileSink{base: basePath, MaxBytes: maxBytes}
	if err := s.open(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *RotatingFileSink) open() error {
	fs, err := NewFileSink(s.base)
	if err != nil {
		return err
	}
	info, err := os.Stat(s.base)
	if err == nil {
		s.size = info.Size()
	}
	s.current = fs
	return nil
}

func (s *RotatingFileSink) rotate() error {
	_ = s.current.Close()
	archive := fmt.Sprintf("%s.%s", s.base, time.Now().Format("20060102T150405"))
	if err := os.Rename(s.base, archive); err != nil {
		return fmt.Errorf("rotating: rename: %w", err)
	}
	s.size = 0
	return s.open()
}

// Write writes the entry, rotating the file if necessary.
func (s *RotatingFileSink) Write(entry logpipe.Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.size >= s.MaxBytes {
		if err := s.rotate(); err != nil {
			return err
		}
	}
	if err := s.current.Write(entry); err != nil {
		return err
	}
	// approximate: each JSON entry is at least a few bytes; stat lazily
	info, err := os.Stat(s.base)
	if err == nil {
		s.size = info.Size()
	}
	return nil
}

// Close closes the current underlying file.
func (s *RotatingFileSink) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.current.Close()
}
