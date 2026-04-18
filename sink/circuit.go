package sink

import (
	"errors"
	"sync"
	"time"

	"github.com/your/logpipe"
)

// ErrCircuitOpen is returned when the circuit breaker is open.
var ErrCircuitOpen = errors.New("circuit breaker open")

type state int

const (
	stateClosed state = iota
	stateOpen
)

// circuitSink wraps a Sink with a circuit breaker that opens after
// consecutive failures and resets after a cooldown period.
type circuitSink struct {
	mu        sync.Mutex
	inner     logpipe.Sink
	maxFails  int
	cooldown  time.Duration
	fails     int
	state     state
	openedAt  time.Time
}

// NewCircuitSink returns a Sink that stops forwarding entries to inner after
// maxFails consecutive errors, reopening after cooldown.
func NewCircuitSink(inner logpipe.Sink, maxFails int, cooldown time.Duration) logpipe.Sink {
	return &circuitSink{
		inner:    inner,
		maxFails: maxFails,
		cooldown: cooldown,
	}
}

func (c *circuitSink) Write(entry logpipe.Entry) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state == stateOpen {
		if time.Since(c.openedAt) < c.cooldown {
			return ErrCircuitOpen
		}
		// attempt reset
		c.state = stateClosed
		c.fails = 0
	}

	err := c.inner.Write(entry)
	if err != nil {
		c.fails++
		if c.fails >= c.maxFails {
			c.state = stateOpen
			c.openedAt = time.Now()
		}
		return err
	}
	c.fails = 0
	return nil
}

func (c *circuitSink) Close() error {
	return c.inner.Close()
}
