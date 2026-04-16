package sink

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/example/logpipe"
)

// ConsoleSink writes JSON-encoded log entries to an io.Writer.
type ConsoleSink struct {
	w       io.Writer
	pretty  bool
}

// NewConsoleSink returns a ConsoleSink writing to w.
// If pretty is true, output is indented.
func NewConsoleSink(w io.Writer, pretty bool) *ConsoleSink {
	if w == nil {
		w = os.Stdout
	}
	return &ConsoleSink{w: w, pretty: pretty}
}

// Write serialises the entry and writes it followed by a newline.
func (c *ConsoleSink) Write(entry logpipe.Entry) error {
	var (b   []byte
		 err error)
	if c.pretty {
		b, err = json.MarshalIndent(entry, "", "  ")
	} else {
		b, err = json.Marshal(entry)
	}
	if err != nil {
		return fmt.Errorf("console sink marshal: %w", err)
	}
	b = append(b, '\n')
	_, err = c.w.Write(b)
	return err
}

// Close is a no-op for the console sink.
func (c *ConsoleSink) Close() error { return nil }
