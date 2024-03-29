/*
 * @Author: cnzf1
 * @Date: 2022-08-26 11:42:04
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-03-28 20:53:28
 * @Description:
 */
package timex

import (
	"math/rand"
	"time"
)

// JitterUp return duration which added factor times the cardinality
// For example for 10s and jitter 1, it will return a time within [10s, 20s])
func JitterUp(duration time.Duration, factor float64) time.Duration {
	if factor <= 0.0 {
		factor = 1.0
	}

	return duration + time.Duration(rand.Float64()*factor*float64(duration))
}

// JitterAround return duration which added a random jitter
//
// This adds or subtracts time from the duration within a given jitter fraction.
// For example for 10s and jitter 0.1, it will return a time within [9s, 11s])
func JitterAround(duration time.Duration, jitter float64) time.Duration {
	multiplier := jitter * (rand.Float64()*2 - 1)
	return time.Duration(float64(duration) * (1 + multiplier))
}

type BackoffManager interface {
	Backoff() Ticker
}

type jitteredBackoffManagerImpl struct {
	timer    Ticker
	duration time.Duration
	jitter   float64
}

// NewJitteredBackoffManager returns a BackoffManager that backoffs with given duration plus given jitter. If the jitter
// is negative, backoff will not be jittered.
func NewJitteredBackoffManager(duration time.Duration, jitter float64) BackoffManager {
	return &jitteredBackoffManagerImpl{
		duration: duration,
		jitter:   jitter,
		timer:    nil,
	}
}

func (j *jitteredBackoffManagerImpl) getNextBackoff() time.Duration {
	jitteredPeriod := j.duration
	if j.jitter > 0.0 {
		jitteredPeriod = JitterUp(j.duration, j.jitter)
	}
	return jitteredPeriod
}

// Backoff implements BackoffManager.Backoff, it returns a timer so caller can block on the timer for jittered backoff.
// The returned timer must be drained before calling Backoff() the second time
func (j *jitteredBackoffManagerImpl) Backoff() Ticker {
	backoff := j.getNextBackoff()
	if j.timer == nil {
		j.timer = NewTicker(backoff)
	} else {
		j.timer.Reset(backoff)
	}
	return j.timer
}

func Until(f func(), period time.Duration, stopCh <-chan struct{}) {
	JitterUntil(f, period, 0.0, true, stopCh)
}

func JitterUntil(f func(), period time.Duration, jitterFactor float64, sliding bool, stopCh <-chan struct{}) {
	BackoffUntil(f, NewJitteredBackoffManager(period, jitterFactor), sliding, stopCh)
}

func BackoffUntil(f func(), backoff BackoffManager, sliding bool, stopCh <-chan struct{}) {
	var t Ticker
	for {
		select {
		case <-stopCh:
			return
		default:
		}

		if !sliding {
			t = backoff.Backoff()
		}

		f()

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
			t.Stop()
			return
		case <-t.Chan():
		}
	}
}
