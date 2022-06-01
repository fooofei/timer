package timer

import (
	"sync/atomic"
	"time"
)

// Timer is a wrapper of time.Timer for correct use of timer.Reset

type Timer struct {
	T      *time.Timer
	active int64
}

func New(d time.Duration) *Timer {
	return &Timer{
		T:      time.NewTimer(d),
		active: 1,
	}
}

// Wait just returns t.T.C for read
// MUST call t.SetUnActive() after read Wait() success
func (t *Timer) Wait() <-chan time.Time {
	return t.T.C
}

// After will restart a new timer with dur
// returns t.Wait() for read
// The restart is safely, auto stop the old timer and keep channel clean
func (t *Timer) After(dur time.Duration) <-chan time.Time {
	t.Reset(dur)
	return t.Wait()
}

// SetUnActive will mark the t.T.C is drained
// please must call it when Wait() returned
func (t *Timer) SetUnActive() {
	atomic.AddInt64(&t.active, -1)
}

func (t *Timer) setActive() {
	atomic.AddInt64(&t.active, 1)
}

// Stop will stop timer safely than time.Timer Stop()
// It will read the dirty data from t.T.C for keep t.T.C clean
func (t *Timer) Stop() bool {
	success := t.T.Stop()
	if !success && atomic.LoadInt64(&t.active) > 1 {
		<-t.T.C
		t.SetUnActive()
	}
	if success {
		// when stop success, no more sendTimer, the T.C is empty,
		// the the un-active status, not read from T.C
		t.SetUnActive()
	}
	return success
}

// Reset will reuse a timer with a new time start
// old timer will be stop safely
func (t *Timer) Reset(dur time.Duration) bool {
	t.Stop()
	active := t.T.Reset(dur)
	t.setActive()
	return active
}

// Go's timer:
// 1、Not stopped timer is in active status, and timer.C maybe pushed sometime
// 2、Go ask us to call timer.Stop() before timer.Reset
// 3、If t.Stop() return true, then timer.C will never be pushed. When t.Stop() return false,
//  timer.C maybe empty or not :
// 		if read from timer.C before timer.Stop(), then timer.C is empty
// 	    if not read from timer.C before timer.Stop(), timer.C maybe empty after timer.Stop() after a while will be write
//            so we cannot use select-default or len to detect timer.C is empty or not, we bind a value `active` to mark it
// 4、Please not use timer in multi go routine

// ref
// https://tonybai.com/2016/12/21/how-to-use-timer-reset-in-golang-correctly/
// https://github.com/golang/go/issues/11513
// https://github.com/golang/go/issues/14038
// another implementation  https://github.com/desertbit/timer
