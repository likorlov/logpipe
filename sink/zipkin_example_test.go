package sink_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/andyyhope/logpipe"
	"github.com/andyyhope/logpipe/sink"
)

func ExampleNewZipkinSink() {
	// Minimal stub server that mimics a Zipkin collector.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	// Create a sink that ships spans to the stub.
	s := sink.NewZipkinSink(
		ts.URL+"/api/v2/spans",
		sink.WithZipkinNameField("op"),
	)
	defer s.Close()

	err := s.Write(logpipe.Entry{
		Fields: logpipe.Fields{
			"op":     "cache.get",
			"key":    "session:abc",
			"hit":    "true",
		},
	})
	fmt.Println(err)
	// Output: <nil>
}
