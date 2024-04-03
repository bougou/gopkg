package wait

import (
	"time"

	"github.com/bougou/gopkg/timeutil"
)

type jitteredBackoffManagerImpl struct {
	clock        timeutil.Clock
	duration     time.Duration
	jitter       float64
	backoffTimer timeutil.Timer
}

// NewJitteredBackoffManager returns a BackoffManager that backoffs with given duration plus given jitter. If the jitter
// is negative, backoff will not be jittered.
func NewJitteredBackoffManager(duration time.Duration, jitter float64, c timeutil.Clock) BackoffManager {
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
func (j *jitteredBackoffManagerImpl) Backoff() timeutil.Timer {
	backoff := j.getNextBackoff()
	if j.backoffTimer == nil {
		j.backoffTimer = j.clock.NewTimer(backoff)
	} else {
		j.backoffTimer.Reset(backoff)
	}
	return j.backoffTimer
}
