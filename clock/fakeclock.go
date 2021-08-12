package clock

import (
	"sync"
	"time"
)

// FakePassiveClock implements PassiveClock, but returns an arbitrary time.
type FakePassiveClock struct {
	lock sync.RWMutex
	time time.Time
}

// FakeClock implements Clock, but returns an arbitrary time.
type FakeClock struct {
	FakePassiveClock

	// waiters are waiting for the fake time to pass their specified time
	waiters []fakeClockWaiter
}

type fakeClockWaiter struct {
	targetTime    time.Time
	stepInterval  time.Duration
	skipIfBlocked bool
	// destChan 通过将 destChan 赋给实际的 Timer 或 Ticker 内部的 chan，这样 Timer/Ticker 就称为了 clock 的 waiter
	destChan chan time.Time
}

// NewFakePassiveClock returns a new FakePassiveClock.
func NewFakePassiveClock(t time.Time) *FakePassiveClock {
	return &FakePassiveClock{
		time: t,
	}
}

// NewFakeClock returns a new FakeClock.
func NewFakeClock(t time.Time) *FakeClock {
	return &FakeClock{
		FakePassiveClock: *NewFakePassiveClock(t),
	}
}

// Now returns f's time.
func (f *FakePassiveClock) Now() time.Time {
	f.lock.RLock()
	defer f.lock.RUnlock()
	return f.time
}

// Since returns time since the time in f.
func (f *FakePassiveClock) Since(ts time.Time) time.Duration {
	f.lock.RLock()
	defer f.lock.RUnlock()
	return f.time.Sub(ts)
}

// SetTime sets the time on the FakePassiveClock.
func (f *FakePassiveClock) SetTime(t time.Time) {
	f.lock.RLock()
	defer f.lock.RUnlock()
	f.time = t
}

// After is the Fake version of time.After(d).
func (f *FakeClock) After(d time.Duration) <-chan time.Time {
	f.lock.Lock()
	defer f.lock.Unlock()
	stopTime := f.time.Add(d)
	ch := make(chan time.Time, 1) // Don't block!
	f.waiters = append(f.waiters, fakeClockWaiter{
		targetTime: stopTime,
		destChan:   ch,
	})
	return ch
}

// NewTimer is the Fake version of time.NewTimer(d).
func (f *FakeClock) NewTimer(d time.Duration) Timer {
	f.lock.Lock()
	defer f.lock.Unlock()
	stopTime := f.time.Add(d)
	ch := make(chan time.Time, 1) // Don't block!
	timer := &fakeTimer{
		fakeClock: f,
		waiter: fakeClockWaiter{
			targetTime: stopTime,
			destChan:   ch,
		},
	}
	f.waiters = append(f.waiters, timer.waiter)
	return timer
}

// NewTicker returns a new Ticker.
func (f *FakeClock) NewTicker(d time.Duration) Ticker {
	f.lock.Lock()
	defer f.lock.Unlock()
	tickTime := f.time.Add(d)
	ch := make(chan time.Time, 1) // Don't block!
	f.waiters = append(f.waiters, fakeClockWaiter{
		targetTime:    tickTime,
		stepInterval:  d,
		skipIfBlocked: true,
		destChan:      ch,
	})

	return &fakeTicker{
		c: ch,
	}
}

// Step moves clock by Duration, notifies anyone that's called After, Tick, or NewTimer
func (f *FakeClock) Step(d time.Duration) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.setTimeLocked(f.time.Add(d))
}

// SetTime sets the time on a FakeClock.
func (f *FakeClock) SetTime(t time.Time) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.setTimeLocked(t)
}

// Actually changes the time and checks any waiters. f must be write-locked.
func (f *FakeClock) setTimeLocked(t time.Time) {
	f.time = t
	newWaiters := make([]fakeClockWaiter, 0, len(f.waiters))
	for i := range f.waiters {
		w := &f.waiters[i]
		if !w.targetTime.After(t) {

			if w.skipIfBlocked {
				select {
				case w.destChan <- t:
				default:
				}
			} else {
				w.destChan <- t
			}

			if w.stepInterval > 0 {
				for !w.targetTime.After(t) {
					w.targetTime = w.targetTime.Add(w.stepInterval)
				}
				newWaiters = append(newWaiters, *w)
			}

		} else {
			newWaiters = append(newWaiters, f.waiters[i])
		}
	}
	f.waiters = newWaiters
}

// HasWaiters returns true if After has been called on f but not yet satisfied (so you can
// write race-free tests).
func (f *FakeClock) HasWaiters() bool {
	f.lock.RLock()
	defer f.lock.RUnlock()
	return len(f.waiters) > 0
}

// Sleep pauses the FakeClock for duration d.
func (f *FakeClock) Sleep(d time.Duration) {
	f.Step(d)
}
