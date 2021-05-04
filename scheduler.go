package concurrent

import (
	"sync"
	"time"
)

// A FixedDelayScheduler executes within a fixed interval. It delays its firing interval by the actual
// callback execution duration, so it is guaranteed that at most only one callback runs at a single point of time
// and that the time between scheduled jobs matches the fixed delay.
type FixedDelayScheduler struct {
	stop   chan bool
	timer  *time.Timer
	mutex  sync.Mutex
	closed bool
}

// NewFixedDelayScheduler allocates a new scheduler and schedules the first interval.
func NewFixedDelayScheduler(delay time.Duration, f func()) *FixedDelayScheduler {
	scheduler := &FixedDelayScheduler{}
	scheduler.stop = make(chan bool)
	scheduler.timer = time.NewTimer(delay)

	go func() {
		defer scheduler.Stop() // close in case of panic or close

		for {
			select {
			// the timer is leaked at most the delay time, but that does not affect
			// self or the callback f, so both can be freed earlier.
			case <-scheduler.timer.C:
			case <-scheduler.stop:
				return
			}

			f()

			scheduler.timer.Reset(delay)
		}
	}()

	return scheduler
}

// Stop deallocates the next scheduled job. If a job is currently running, it is not interrupted.
// Stop is idempotent and thread safe.
func (s *FixedDelayScheduler) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// avoid close of closed channel
	if s.closed {
		return
	}

	s.closed = true

	close(s.stop)
	s.timer.Stop()
}
