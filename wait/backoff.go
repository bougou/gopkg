package wait

import (
	"math"
	"time"

	"github.com/bougou/gopkg/clock"
)

// Backoff 退避
// 在一个操作发生冲突后，等待一定时间后再进行
// 退避算法有几种：
// 按照固定时间间隔重试:
// 按照指数时间间隔重试:
type Backoff struct {
	// The initial duration.
	Duration time.Duration
	// Duration is multiplied by factor each iteration, if factor is not zero
	// and the limits imposed by Steps and Cap have not been reached.
	// Should not be negative.
	// The jitter does not contribute to the updates to the duration parameter.
	Factor float64
	// The sleep at each iteration is the duration plus an additional
	// amount chosen uniformly at random from the interval between
	// zero and `jitter*duration`.
	Jitter float64
	// The remaining number of iterations in which the duration
	// parameter may change (but progress can be stopped earlier by
	// hitting the cap). If not positive, the duration is not
	// changed. Used for exponential backoff in combination with
	// Factor and Cap.
	Steps int
	// A limit on revised values of the duration parameter. If a
	// multiplication by the factor parameter would make the duration
	// exceed the cap then the duration is set to the cap and the
	// steps parameter is set to zero.
	Cap time.Duration
}

// Step (1) returns an amount of time to sleep determined by the
// original Duration and Jitter and (2) mutates the provided Backoff
// to update its Steps and Duration.
func (b *Backoff) Step() time.Duration {
	if b.Steps < 1 {
		if b.Jitter > 0 {
			return Jitter(b.Duration, b.Jitter)
		}
		return b.Duration
	}

	b.Steps--

	duration := b.Duration

	// calculate the next step
	if b.Factor != 0 {
		b.Duration = time.Duration(float64(b.Duration) * b.Factor)
		if b.Cap > 0 && b.Duration < b.Cap {
			b.Duration = b.Cap
			b.Steps = 0
		}
	}

	if b.Jitter > 0 {
		duration = Jitter(duration, b.Jitter)
	}

	return duration
}

// BackoffManager manages backoff with a particular scheme based on its underlying implementation. It provides
// an interface to return a timer for backoff, and caller shall backoff until Timer.C() drains. If the second Backoff()
// is called before the timer from the first Backoff() call finishes, the first timer will NOT be drained and result in
// undetermined behavior.
// The BackoffManager is supposed to be called in a single-threaded environment.
type BackoffManager interface {
	Backoff() clock.Timer
}

type exponentialBackoffManagerImpl struct {
	backoff              *Backoff
	backoffTimer         clock.Timer
	lastBackoffStart     time.Time
	initialBackoff       time.Duration
	backoffResetDuration time.Duration
	clock                clock.Clock
}

// NewExponentialBackoffManager returns a manager for managing exponential backoff. Each backoff is jittered and
// backoff will not exceed the given max. If the backoff is not called within resetDuration, the backoff is reset.
// This backoff manager is used to reduce load during upstream unhealthiness.
func NewExponentialBackoffManager(initBackoff, maxBackoff, resetDuration time.Duration, backoffFactor, jitter float64, c clock.Clock) BackoffManager {
	return &exponentialBackoffManagerImpl{
		backoff: &Backoff{
			Duration: initBackoff,
			Factor:   backoffFactor,
			Jitter:   jitter,

			// the current impl of wait.Backoff returns Backoff.Duration once steps are used up, which is not
			// what we ideally need here, we set it to max int and assume we will never use up the steps
			Steps: math.MaxInt32,
			Cap:   maxBackoff,
		},
		backoffTimer:         nil,
		initialBackoff:       initBackoff,
		lastBackoffStart:     c.Now(),
		backoffResetDuration: resetDuration,
		clock:                c,
	}
}

func (b *exponentialBackoffManagerImpl) getNextBackoff() time.Duration {
	if b.clock.Now().Sub(b.lastBackoffStart) > b.backoffResetDuration {
		b.backoff.Steps = math.MaxInt32
		b.backoff.Duration = b.initialBackoff
	}
	b.lastBackoffStart = b.clock.Now()
	return b.backoff.Step()
}

// Backoff implements BackoffManager.Backoff, it returns a timer so caller can block on the timer for exponential backoff.
// The returned timer must be drained before calling Backoff() the second time
func (b *exponentialBackoffManagerImpl) Backoff() clock.Timer {
	if b.backoffTimer == nil {
		b.backoffTimer = b.clock.NewTimer(b.getNextBackoff())
	} else {
		b.backoffTimer.Reset(b.getNextBackoff())
	}
	return b.backoffTimer
}

type jitteredBackoffManagerImpl struct {
	clock        clock.Clock
	duration     time.Duration
	jitter       float64
	backoffTimer clock.Timer
}

// NewJitteredBackoffManager returns a BackoffManager that backoffs with given duration plus given jitter. If the jitter
// is negative, backoff will not be jittered.
func NewJitteredBackoffManager(duration time.Duration, jitter float64, c clock.Clock) BackoffManager {
	return &jitteredBackoffManagerImpl{
		clock:        c,
		duration:     duration,
		jitter:       jitter,
		backoffTimer: nil,
	}
}

func (j *jitteredBackoffManagerImpl) getNextBackoff() time.Duration {
	jitteredPeriod := j.duration
	if j.jitter > 0.0 {
		jitteredPeriod = Jitter(j.duration, j.jitter)
	}
	return jitteredPeriod
}

// Backoff implements BackoffManager.Backoff, it returns a timer so caller can block on the timer for jittered backoff.
// The returned timer must be drained before calling Backoff() the second time
func (j *jitteredBackoffManagerImpl) Backoff() clock.Timer {
	backoff := j.getNextBackoff()
	if j.backoffTimer == nil {
		j.backoffTimer = j.clock.NewTimer(backoff)
	} else {
		j.backoffTimer.Reset(backoff)
	}
	return j.backoffTimer
}

// ExponentialBackoff repeats a condition check with exponential backoff.
//
// It repeatedly checks the condition and then sleeps, using `backoff.Step()`
// to determine the length of the sleep and adjust Duration and Steps.
// Stops and returns as soon as:
// 1. the condition check returns true or an error,
// 2. `backoff.Steps` checks of the condition have been done, or
// 3. a sleep truncated by the cap on duration has been completed.
// In case (1) the returned error is what the condition function returned.
// In all other cases, ErrWaitTimeout is returned.
func ExponentialBackoff(backoff Backoff, condition ConditionFunc) error {
	for backoff.Steps > 0 {
		if ok, err := runConditionWithCrashProtection(condition); err != nil || ok {
			return err
		}
		if backoff.Steps == 1 {
			break
		}
		time.Sleep(backoff.Step())
	}
	return ErrWaitTimeout
}
