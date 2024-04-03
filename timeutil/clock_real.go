package timeutil

import "time"

// RealClock really calls time.Now()
type RealClock struct{}

// Now returns the current time.
func (RealClock) Now() time.Time {
	return time.Now()
}

// Since returns time since the specified timestamp.
func (RealClock) Since(ts time.Time) time.Duration {
	return time.Since(ts)
}

// After is the same as time.After(d).
func (RealClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

// NewTimer returns a new Timer.
func (RealClock) NewTimer(d time.Duration) Timer {
	return &realTimer{
		timer: time.NewTimer(d),
	}
}

// NewTicker returns a new Ticker.
func (RealClock) NewTicker(d time.Duration) Ticker {
	return &realTicker{
		ticker: time.NewTicker(d),
	}
}

// Sleep pauses the RealClock for duration d.
func (RealClock) Sleep(d time.Duration) {
	time.Sleep(d)
}
