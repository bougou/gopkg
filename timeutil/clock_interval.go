package timeutil

import "time"

// IntervalClock implements Clock, but each invocation of Now steps the clock forward the specified duration
// 间隔式时钟，（有一个起始时刻和一个间隔时段）
// 每次调用 Now，返回的时间是：上次记录的时刻 + 间隔时段；这个返回的时间会被记录下来。
type IntervalClock struct {
	timee    time.Time
	duration time.Duration
}

// IntervalClock implements Clock interface
var _ Clock = (*IntervalClock)(nil)

func NewIntervalClock(timee time.Time, duration time.Duration) IntervalClock {
	return IntervalClock{
		timee:    timee,
		duration: duration,
	}
}

// Now returns i's time
func (i *IntervalClock) Now() time.Time {
	i.timee = i.timee.Add(i.duration)
	return i.timee
}

// Since returns time since the time in i.
func (i *IntervalClock) Since(ts time.Time) time.Duration {
	return i.timee.Sub(ts)
}

// After is currently unimplemented, will panic.
// TODO: make interval clock use FakeClock so this can be implemented.
func (*IntervalClock) After(d time.Duration) <-chan time.Time {
	panic("IntervalClock doesn't implement After")
}

// NewTimer is currently unimplemented, will panic.
// TODO: make interval clock use FakeClock so this can be implemented.
func (*IntervalClock) NewTimer(d time.Duration) Timer {
	panic("IntervalClock doesn't implement NewTimer")
}

// NewTicker is currently unimplemented, will panic.
// TODO: make interval clock use FakeClock so this can be implemented.
func (*IntervalClock) NewTicker(d time.Duration) Ticker {
	panic("IntervalClock doesn't implement NewTicker")
}

// Sleep is currently unimplemented; will panic.
func (*IntervalClock) Sleep(d time.Duration) {
	panic("IntervalClock doesn't implement Sleep")
}
