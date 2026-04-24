package sink_test

import (
	"errors"
	"testing"

	"github.com/logpipe/logpipe"
	"github.com/logpipe/logpipe/sink"
)

func TestStashSink_StoresAndRetrieves(t *testing.T) {
	col := &collectSink{}
	ss := sink.NewStashSink(col, func(e logpipe.Entry) string {
		if v, ok := e["id"].(string); ok {
			return v
		}
		return ""
	})

	e := logpipe.Entry{"id": "req-1", "msg": "hello"}
	if err := ss.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, ok := ss.Get("req-1")
	if !ok {
		t.Fatal("expected entry to be stashed")
	}
	if got["msg"] != "hello" {
		t.Fatalf("unexpected msg: %v", got["msg"])
	}
	if len(col.entries) != 1 {
		t.Fatalf("expected 1 forwarded entry, got %d", len(col.entries))
	}
}

func TestStashSink_OverwritesSameKey(t *testing.T) {
	ss := sink.NewStashSink(&collectSink{}, func(e logpipe.Entry) string {
		return "key"
	})
	_ = ss.Write(logpipe.Entry{"v": 1})
	_ = ss.Write(logpipe.Entry{"v": 2})

	got, _ := ss.Get("key")
	if got["v"] != 2 {
		t.Fatalf("expected overwritten value 2, got %v", got["v"])
	}
	if ss.Len() != 1 {
		t.Fatalf("expected len 1, got %d", ss.Len())
	}
}

func TestStashSink_EmptyKeyNotStashed(t *testing.T) {
	ss := sink.NewStashSink(&collectSink{}, func(e logpipe.Entry) string {
		return ""
	})
	_ = ss.Write(logpipe.Entry{"msg": "no key"})
	if ss.Len() != 0 {
		t.Fatalf("expected nothing stashed, got %d", ss.Len())
	}
}

func TestStashSink_Delete(t *testing.T) {
	ss := sink.NewStashSink(&collectSink{}, func(e logpipe.Entry) string { return "k" })
	_ = ss.Write(logpipe.Entry{"x": 1})
	ss.Delete("k")
	_, ok := ss.Get("k")
	if ok {
		t.Fatal("expected entry to be deleted")
	}
}

func TestStashSink_PropagatesError(t *testing.T) {
	sentinel := errors.New("sink error")
	ss := sink.NewStashSink(errorSink(sentinel), func(e logpipe.Entry) string { return "k" })
	if err := ss.Write(logpipe.Entry{}); !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestStashSink_Close(t *testing.T) {
	col := &collectSink{}
	ss := sink.NewStashSink(col, func(e logpipe.Entry) string { return "k" })
	_ = ss.Write(logpipe.Entry{"a": 1})
	if err := ss.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
	if ss.Len() != 0 {
		t.Fatal("expected stash cleared after close")
	}
}

func TestStashSink_PanicOnNilKeyFn(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil keyFn")
		}
	}()
	sink.NewStashSink(&collectSink{}, nil)
}
