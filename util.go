package main

import (
	"sync"
	"time"
)

type rateState struct {
	mu        sync.RWMutex  // mutex
	lastFlush time.Time     // last update time
	period    time.Duration // observe time period
	tickCount int           // count
	rps       float64       // rps
}

func (pRate *rateState) ticRate(i int) {
	pRate.mu.Lock()
	pRate.tickCount += i
	now, tps := time.Now(), 0.
	if now.Sub(pRate.lastFlush) >= pRate.period {
		if pRate.tickCount > 0 {
			tps = float64(pRate.tickCount) / now.Sub(pRate.lastFlush).Seconds()
		}
		pRate.tickCount, pRate.lastFlush, pRate.rps = 0, now, tps
	}
	pRate.mu.Unlock()
}

func (pRate *rateState) getRate() float64 {
	pRate.ticRate(0)
	t := 0.
	pRate.mu.RLock()
	t = pRate.rps
	pRate.mu.RUnlock()
	return t
}
