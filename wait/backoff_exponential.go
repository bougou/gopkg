package wait

import (
	"math"
	"time"

	"github.com/bougou/gopkg/timeutil"
)

type exponentialBackoffManagerImpl struct {
	backoff              *Backoff
	backoffTimer         timeutil.Timer
	lastBackoffStart     time.Time
	initialBackoff       time.Duration
	backoffResetDuration time.Duration
	clock                timeutil.Clock
}

// NewExponentialBackoffManager returns a manager for managing exponential backoff. Each backoff is jittered and
// backoff will not exceed the given max. If the backoff is not called within resetDuration, the backoff is reset.
// This backoff manager is used to reduce load during upstream unhealthiness.
func NewExponentialBackoffManager(initBackoff, maxBackoff, resetDuration time.Duration, backoffFactor, jitter float64, c timeutil.Clock) BackoffManager {
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
func (b *exponentialBackoffManagerImpl) Backoff() timeutil.Timer {
	if b.backoffTimer == nil {
		b.backoffTimer = b.clock.NewTimer(b.getNextBackoff())
	} else {
		b.backoffTimer.Reset(b.getNextBackoff())
	}
	return b.backoffTimer
}
