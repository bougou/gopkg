package timeutil

import (
	"time"
)

// PassiveClock allows for injecting fake or real clocks into code
// that needs to read the current time but does not support scheduling
// activity in the future.
type PassiveClock interface {
	Now() time.Time
	Since(time.Time) time.Duration
}

// Clock allows for injecting fake or real clocks into code that
// needs to do arbitrary things based on time.
type Clock interface {
	PassiveClock
	After(time.Duration) <-chan time.Time
	NewTimer(time.Duration) Timer
	Sleep(time.Duration)
	NewTicker(time.Duration) Ticker
}
