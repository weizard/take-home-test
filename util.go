package main

import (
	"fmt"
	"os"
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

var f *os.File

func errorLog(title string, content string) {
	logFileName := "log"
	var err error
	if f == nil {
		if f, err = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			fmt.Printf("[errorLog] %s\n", err.Error())
		}
	}
	f.WriteString(time.Now().Format("2006-01-02 03:04:05 -0700") + " - [" + title + "] " + content + "\n")
}
