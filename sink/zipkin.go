package sink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/andyyhope/logpipe"
)

// zipkinSpan is the subset of Zipkin v2 span fields we populate.
type zipkinSpan struct {
	Name      string            `json:"name"`
	Timestamp int64             `json:"timestamp"` // microseconds since epoch
	Tags      map[string]string `json:"tags,omitempty"`
}

// zipkinSink ships log entries to a Zipkin-compatible HTTP endpoint as spans.
type zipkinSink struct {
	endpoint string
	client   *http.Client
	field    string
}

// ZipkinOption configures a zipkinSink.
type ZipkinOption func(*zipkinSink)

// WithZipkinNameField sets the entry field used as the span name (default: "message").
func WithZipkinNameField(field string) ZipkinOption {
	return func(z *zipkinSink) { z.field = field }
}

// WithZipkinHTTPClient replaces the default HTTP client.
func WithZipkinHTTPClient(c *http.Client) ZipkinOption {
	return func(z *zipkinSink) { z.client = c }
}

// NewZipkinSink returns a Sink that forwards each log entry to a Zipkin
// HTTP/JSON endpoint (e.g. http://localhost:9411/api/v2/spans).
func NewZipkinSink(endpoint string, opts ...ZipkinOption) logpipe.Sink {
	z := &zipkinSink{
		endpoint: endpoint,
		client:   &http.Client{Timeout: 5 * time.Second},
		field:    "message",
	}
	for _, o := range opts {
		o(z)
	}
	return z
}

func (z *zipkinSink) Write(e logpipe.Entry) error {
	name, _ := e.Fields[z.field].(string)
	if name == "" {
		name = "log"
	}
	tags := make(map[string]string, len(e.Fields))
	for k, v := range e.Fields {
		tags[k] = fmt.Sprintf("%v", v)
	}
	span := zipkinSpan{
		Name:      name,
		Timestamp: time.Now().UnixMicro(),
		Tags:      tags,
	}
	buf, err := json.Marshal([]zipkinSpan{span})
	if err != nil {
		return fmt.Errorf("zipkin: marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, z.endpoint, bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("zipkin: request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := z.client.Do(req)
	if err != nil {
		return fmt.Errorf("zipkin: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("zipkin: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func (z *zipkinSink) Close() error { return nil }
