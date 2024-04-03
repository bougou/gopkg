package wait

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/bougou/gopkg/runtime"
	"github.com/bougou/gopkg/timeutil"
)

// For any test of the style:
//
//	...
//	<- time.After(timeout):
//	   t.Errorf("Timed out")
//
// The value for timeout should effectively be "forever." Obviously we don't want our tests to truly lock up forever, but 30s
// is long enough that it is effectively forever for the things that can slow down a run on a heavily contended machine
// (GC, seeks, etc), but not so long as to make a developer ctrl-c a test run if they do happen to break that test.
var ForeverTestTimeout = time.Second * 30

// NeverStop may be passed to Until to make it never stop.
var NeverStop <-chan struct{} = make(chan struct{})

// ErrWaitTimeout is returned when the condition exited without success.
var ErrWaitTimeout = errors.New("timed out waiting for the condition")

// ConditionFunc returns true if the condition is satisfied, or an error
// if the loop should be aborted.
type ConditionFunc func() (done bool, err error)

// runConditionWithCrashProtection runs a ConditionFunc with crash protection
func runConditionWithCrashProtection(condition ConditionFunc) (bool, error) {
	defer runtime.HandleCrash()
	return condition()
}

// contextForChannel derives a child context from a parent channel.
//
// The derived context's Done channel is closed when the returned cancel function
// is called or when the parent channel is closed, whichever happens first.
//
// Note the caller must *always* call the CancelFunc, otherwise resources may be leaked.
func contextForChannel(parentCh <-chan struct{}) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		select {
		case <-parentCh:
			cancel()
		case <-ctx.Done():
		}
	}()
	return ctx, cancel
}

// Jitter returns a time.Duration between duration and duration + maxFactor *
// duration.
//
// This allows clients to avoid converging on periodic behavior. If maxFactor
// is 0.0, a suggested default value will be chosen.
//
// Jitter 返回一个指定间隔范围内的间隔
func Jitter(duration time.Duration, maxFactor float64) time.Duration {
	if maxFactor <= 0.0 {
		maxFactor = 1.0
	}
	wait := duration + time.Duration(rand.Float64()*maxFactor*float64(duration))
	return wait
}

// Forever calls f every period for ever.
//
// Forever is syntactic sugar on top of Until.
func Forever(f func(), period time.Duration) {
	Until(f, period, NeverStop)
}

// Until loops until stop channel is closed, running f every period.
//
// Until is syntactic sugar on top of JitterUntil with zero jitter factor and
// with sliding = true (which means the timer for period starts after the f
// completes).
func Until(f func(), period time.Duration, stopCh <-chan struct{}) {
	JitterUntil(f, period, 0.0, true, stopCh)
}

// UntilWithContext loops until context is done, running f every period.
//
// UntilWithContext is syntactic sugar on top of JitterUntilWithContext
// with zero jitter factor and with sliding = true (which means the timer
// for period starts after the f completes).
func UntilWithContext(ctx context.Context, f func(context.Context), period time.Duration) {
	JitterUntilWithContext(ctx, f, period, 0.0, true)
}

// NonSlidingUntil loops until stop channel is closed, running f every
// period.
//
// NonSlidingUntil is syntactic sugar on top of JitterUntil with zero jitter
// factor, with sliding = false (meaning the timer for period starts at the same
// time as the function starts).
func NonSlidingUntil(f func(), period time.Duration, stopCh <-chan struct{}) {
	JitterUntil(f, period, 0.0, false, stopCh)
}

// NonSlidingUntilWithContext loops until context is done, running f every
// period.
//
// NonSlidingUntilWithContext is syntactic sugar on top of JitterUntilWithContext
// with zero jitter factor, with sliding = false (meaning the timer for period
// starts at the same time as the function starts).
func NonSlidingUntilWithContext(ctx context.Context, f func(context.Context), period time.Duration) {
	JitterUntilWithContext(ctx, f, period, 0.0, false)
}

// JitterUntil loops until stop channel is closed, running f every period.
//
// If jitterFactor is positive, the period is jittered before every run of f.
// If jitterFactor is not positive, the period is unchanged and not jittered.
//
// If sliding is true, the period is computed after f runs. If it is false then
// period includes the runtime for f.
//
// Close stopCh to stop. f may not be invoked if stop channel is already
// closed. Pass NeverStop to if you don't want it stop.
func JitterUntil(f func(), period time.Duration, jitterFactor float64, sliding bool, stopCh <-chan struct{}) {
	BackoffUntil(f, NewJitteredBackoffManager(period, jitterFactor, &timeutil.RealClock{}), sliding, stopCh)
}

// BackoffUntil loops until stop channel is closed, run f every duration given by BackoffManager.
//
// If sliding is true, the period is computed after f runs. If it is false then
// period includes the runtime for f.
func BackoffUntil(f func(), backoff BackoffManager, sliding bool, stopCh <-chan struct{}) {
	var t timeutil.Timer
	for {
		select {
		case <-stopCh:
			return
		default:
		}

		if !sliding {
			t = backoff.Backoff()
		}

		func() {
			defer runtime.HandleCrash()
			f()
		}()

		if sliding {
			t = backoff.Backoff()
		}

		// NOTE: b/c there is no priority selection in golang
		// it is possible for this to race, meaning we could
		// trigger t.C and stopCh, and t.C select falls through.
		// In order to mitigate we re-check stopCh at the beginning
		// of every loop to prevent extra executions of f().
		select {
		case <-stopCh:
			return
		case <-t.C():
		}
	}
}

// JitterUntilWithContext loops until context is done, running f every period.
//
// If jitterFactor is positive, the period is jittered before every run of f.
// If jitterFactor is not positive, the period is unchanged and not jittered.
//
// If sliding is true, the period is computed after f runs. If it is false then
// period includes the runtime for f.
//
// Cancel context to stop. f may not be invoked if context is already expired.
func JitterUntilWithContext(ctx context.Context, f func(context.Context), period time.Duration, jitterFactor float64, sliding bool) {
	JitterUntil(func() { f(ctx) }, period, jitterFactor, sliding, ctx.Done())
}
