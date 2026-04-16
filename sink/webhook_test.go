package sink_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourorg/logpipe"
	"github.com/yourorg/logpipe/sink"
)

func TestWebhookSink_Write(t *testing.T) {
	var received logpipe.Entry
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	w := sink.NewWebhookSink(ts.URL)
	e := makeEntry(logpipe.INFO, "hello webhook")
	if err := w.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Message != "hello webhook" {
		t.Errorf("expected message 'hello webhook', got %q", received.Message)
	}
}

func TestWebhookSink_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	w := sink.NewWebhookSink(ts.URL)
	err := w.Write(makeEntry(logpipe.ERROR, "fail"))
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestWebhookSink_Timeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	w := sink.NewWebhookSink(ts.URL, sink.WithTimeout(50*time.Millisecond))
	err := w.Write(makeEntry(logpipe.WARN, "slow"))
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestWebhookSink_Close(t *testing.T) {
	w := sink.NewWebhookSink("http://localhost")
	if err := w.Close(); err != nil {
		t.Errorf("Close should return nil, got %v", err)
	}
}
