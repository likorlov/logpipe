package sink_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andyyhope/logpipe"
	"github.com/andyyhope/logpipe/sink"
)

func TestZipkinSink_Write(t *testing.T) {
	var received []map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	s := sink.NewZipkinSink(ts.URL)
	e := logpipe.Entry{Fields: logpipe.Fields{"message": "hello", "env": "test"}}
	if err := s.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received) != 1 {
		t.Fatalf("expected 1 span, got %d", len(received))
	}
	if received[0]["name"] != "hello" {
		t.Errorf("expected name=hello, got %v", received[0]["name"])
	}
}

func TestZipkinSink_CustomNameField(t *testing.T) {
	var received []map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	s := sink.NewZipkinSink(ts.URL, sink.WithZipkinNameField("op"))
	e := logpipe.Entry{Fields: logpipe.Fields{"op": "db.query"}}
	if err := s.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received[0]["name"] != "db.query" {
		t.Errorf("expected name=db.query, got %v", received[0]["name"])
	}
}

func TestZipkinSink_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	s := sink.NewZipkinSink(ts.URL)
	err := s.Write(logpipe.Entry{Fields: logpipe.Fields{"message": "fail"}})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestZipkinSink_FallbackName(t *testing.T) {
	var received []map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	s := sink.NewZipkinSink(ts.URL)
	// No "message" field — should fall back to "log"
	e := logpipe.Entry{Fields: logpipe.Fields{"level": "info"}}
	if err := s.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received[0]["name"] != "log" {
		t.Errorf("expected fallback name=log, got %v", received[0]["name"])
	}
}

func TestZipkinSink_Close(t *testing.T) {
	s := sink.NewZipkinSink("http://localhost:9411")
	if err := s.Close(); err != nil {
		t.Errorf("Close() returned unexpected error: %v", err)
	}
}
