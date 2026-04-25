package sink

import (
	"sync"
	"time"

	"github.com/logpipe/logpipe"
)

// HedgeSink writes an entry to the primary sink and, if it does not complete
// within the hedge delay, concurrently fires the same entry at the secondary
// sink. The first successful write wins; the other result is discarded.
//
// This is useful for reducing tail latency when a remote sink occasionally
// stalls.
type hedgeSink struct {
	primary   logpipe.Sink
	secondary logpipe.Sink
	delay     time.Duration
}

// NewHedgeSink returns a Sink that hedges writes against a secondary sink
// after delay has elapsed without a response from the primary.
func NewHedgeSink(primary, secondary logpipe.Sink, delay time.Duration) logpipe.Sink {
	if delay <= 0 {
		panic("sink: NewHedgeSink delay must be positive")
	}
	return &hedgeSink{primary: primary, secondary: secondary, delay: delay}
}

func (h *hedgeSink) Write(entry logpipe.Entry) error {
	type result struct {
		err error
	}

	primaryCh := make(chan result, 1)
	go func() {
		primaryCh <- result{h.primary.Write(entry)}
	}()

	select {
	case res := <-primaryCh:
		return res.err
	case <-time.After(h.delay):
		// Primary is slow; launch hedge request.
	}

	secondaryCh := make(chan result, 1)
	go func() {
		secondaryCh <- result{h.secondary.Write(entry)}
	}()

	var (
		once sync.Once
		win  error
	)
	setWin := func(err error) { once.Do(func() { win = err }) }

	for i := 0; i < 2; i++ {
		select {
		case res := <-primaryCh:
			if res.err == nil {
				return nil
			}
			setWin(res.err)
		case res := <-secondaryCh:
			if res.err == nil {
				return nil
			}
			setWin(res.err)
		}
	}
	return win
}

func (h *hedgeSink) Close() error {
	pe := h.primary.Close()
	se := h.secondary.Close()
	if pe != nil {
		return pe
	}
	return se
}
