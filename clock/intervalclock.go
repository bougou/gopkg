package clock

import "time"

// IntervalClock implements Clock, but each invocation of Now steps the clock forward the specified duration
// 间隔式时钟，（有一个起始时刻和一个间隔时段）
// 每次调用 Now 时，返回的时间是：上次记录的时刻 + 间隔时段
type IntervalClock struct {
	Time     time.Time
	Duration time.Duration
}

// Now returns i's time
func (i *IntervalClock) Now() time.Time {
	i.Time = i.Time.Add(i.Duration)
	return i.Time
}

// Since returns time since the time in i.
func (i *IntervalClock) Since(ts time.Time) time.Duration {
	return i.Time.Sub(ts)
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
