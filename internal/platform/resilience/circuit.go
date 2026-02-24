package resilience

import (
	"sync"
	"time"

	"github.com/boris989/fulcrum/internal/observability/metrics"
)

type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

type Circuit struct {
	mu sync.Mutex

	state       State
	failures    int
	maxFailures int

	openUntil time.Time
	timeout   time.Duration
}

func New(maxFailures int, timeout time.Duration) *Circuit {
	return &Circuit{
		state:       Closed,
		maxFailures: maxFailures,
		timeout:     timeout,
	}
}

func (c *Circuit) Allow() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state == Open {
		if time.Now().After(c.openUntil) {
			c.state = HalfOpen
			return true
		}

		return false
	}

	return true
}

func (c *Circuit) OnSuccess() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.failures = 0
	c.state = Closed
}

func (c *Circuit) OnFailure() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.failures++

	if c.failures >= c.maxFailures {
		c.state = Open
		c.openUntil = time.Now().Add(c.timeout)
		metrics.KafkaCircuitOpen.Inc()
	}
}
