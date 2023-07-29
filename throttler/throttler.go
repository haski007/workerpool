package throttler

import (
	"sync"
	"time"
)

type Throttler struct {
	invocations int64
	duration    time.Duration
	m           sync.Mutex
	wg          sync.WaitGroup
	max         int64
	counter     int64
}

func New(
	dur time.Duration,
	limit int64,
) *Throttler {
	return &Throttler{
		invocations: 0,
		duration:    dur,
		m:           sync.Mutex{},
		wg:          sync.WaitGroup{},
		max:         limit,
		counter:     0,
	}
}

func (t *Throttler) Throttle() {
	t.m.Lock()
	defer t.m.Unlock()
	t.invocations += 1
	if t.invocations >= t.max {
		<-time.After(t.duration) // This block until the time period expires
		t.invocations = 0
	}
}

func (t *Throttler) Lock() {
	t.m.Lock()
}

func (t *Throttler) Unlock() {
	t.m.Unlock()
}
