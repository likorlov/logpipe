package sink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourorg/logpipe"
)

// WebhookSink sends log entries as JSON POST requests to an HTTP endpoint.
type WebhookSink struct {
	url    string
	client *http.Client
}

// WebhookOption configures a WebhookSink.
type WebhookOption func(*WebhookSink)

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) WebhookOption {
	return func(w *WebhookSink) {
		w.client.Timeout = d
	}
}

// NewWebhookSink creates a sink that POSTs each log entry to the given URL.
func NewWebhookSink(url string, opts ...WebhookOption) *WebhookSink {
	w := &WebhookSink{
		url:    url,
		client: &http.Client{Timeout: 5 * time.Second},
	}
	for _, o := range opts {
		o(w)
	}
	return w
}

// Write encodes the entry as JSON and POSTs it to the configured URL.
func (w *WebhookSink) Write(e logpipe.Entry) error {
	body, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("webhook: marshal: %w", err)
	}
	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// Close is a no-op for WebhookSink.
func (w *WebhookSink) Close() error { return nil }
