package main

import (
	"math"
	"testing"
	"time"
)

func TestRateCounter(t *testing.T) {
	var testWatcher rateState
	testWatcher = rateState{
		lastFlush: time.Now(),
		period:    1 * time.Second,
	}
	testWatcher.ticRate(5)
	time.Sleep(5 * time.Second)
	if math.Floor(testWatcher.getRate()+0.5) != 1. {
		t.Errorf("rate count failed.")
	}
}
