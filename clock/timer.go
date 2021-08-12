package clock

import "time"

// Timer allows for injecting fake or real timers into code that needs
// to do arbitrary things based on time.
// 计时器、定时器
// 1. 周期性发送时间
// 2. 会超时
// 3. 可以停止或重置
type Timer interface {
	C() <-chan time.Time
	Stop() bool
	Reset(d time.Duration) bool
}

// realTimer is backedn by an actual time.Timer.
type realTimer struct {
	timer *time.Timer
}

// C returns the underlying timer's channel.
func (r *realTimer) C() <-chan time.Time {
	return r.timer.C
}

// Stop calls Stop() on the underlying timer.
func (r *realTimer) Stop() bool {
	return r.timer.Stop()
}

// Reset calls Reset() on the underlying timer.
func (r *realTimer) Reset(d time.Duration) bool {
	return r.timer.Reset(d)
}

// fakeTimer implements Timer based on a FakeClock.
type fakeTimer struct {
	fakeClock *FakeClock
	waiter    fakeClockWaiter
}

// C returns the channel that notifies when this timer has fired.
func (f *fakeTimer) C() <-chan time.Time {
	return f.waiter.destChan
}

// Stop conditionally stops the timer. If the timer has neither fired
// nor been stopped then this call stops the timer and returns true,
// otherwise this call returns false.  This is like time.Timer::Stop.
// 有条件地 Stop：如果没有超时或没有被 stopped，则 Stop 停止该定时器，返回 true
func (f *fakeTimer) Stop() bool {
	f.fakeClock.lock.Lock()
	defer f.fakeClock.lock.Unlock()
	// The timer has already fired or been stopped, unless it is found
	// among the clock's waiters.
	// 遍历 该定时器依赖的 clock 中维护的 waiters, 如果
	stopped := false
	oldWaiters := f.fakeClock.waiters
	newWaiters := make([]fakeClockWaiter, 0, len(oldWaiters))
	seekChan := f.waiter.destChan
	for i := range oldWaiters {
		// Identify the timer's by the identity of the
		// destination channel, nothing else is necessarily unique and
		// constant since the timer's creation.
		if oldWaiters[i].destChan == seekChan {
			// 说明在 oldWaiters 定位到了改定时器，意味着该定时器还没有超时或被 stop
			// 则模拟 stop
			stopped = true
		} else {
			newWaiters = append(newWaiters, oldWaiters[i])
		}
	}

	f.fakeClock.waiters = newWaiters
	return stopped
}

// Reset conditionally updates the firing time of the timer. If the
// timer has neither fired nor been stopped then this call resets the
// timer to the fake clock's "now" + d and returns true, otherwise
// this call returns false.  This is like time.Timer::Reset.
func (f *fakeTimer) Reset(d time.Duration) bool {
	f.fakeClock.lock.Lock()
	defer f.fakeClock.lock.Unlock()
	waiters := f.fakeClock.waiters
	seekChan := f.waiter.destChan
	for i := range waiters {
		if waiters[i].destChan == seekChan {
			waiters[i].targetTime = f.fakeClock.time.Add(d)
			return true
		}
	}
	return false
}
