package timer

import (
	"sync/atomic"
	"time"
)

// Timer is a wrapper of time.Timer for correct use of timer.Reset

type Timer struct {
	t      *time.Timer
	active atomic.Int32
}

// New will create a new timer with d for trigger
// You can use as
//
//	    t := New(time.Second)
//		defer t.Stop()
//
//		select {
//		case <-t.After(2 * time.Second):
//			t.SetUnActive()
//		    // do something
//		}
func New(d time.Duration) *Timer {
	var t = &Timer{
		t: time.NewTimer(d),
	}
	t.setActive()
	return t
}

// After is not recommend use
// func After(d time.Duration) <-chan time.Time {
//	return New(d).Wait()
// }

// Wait just returns t.t.C for read
// MUST call t.SetUnActive() after read Wait() success
func (t *Timer) Wait() <-chan time.Time {
	return t.t.C
}

// After will restart a new timer with dur
// returns t.Wait() for read
// The restart is safely, auto stop the old timer and keep channel clean
func (t *Timer) After(dur time.Duration) <-chan time.Time {
	t.Reset(dur)
	return t.Wait()
}

// SetUnActive will mark the t.t.C is drained
// please must call it when Wait() returned
func (t *Timer) SetUnActive() {
	t.active.Add(-1)
}

func (t *Timer) setActive() {
	t.active.Add(1)
}

// Stop will stop timer safely than time.Timer Stop()
// It will read the dirty data from t.t.C for keep t.t.C clean
func (t *Timer) Stop() bool {
	success := t.t.Stop()
	if !success && t.active.Load() > 1 {
		<-t.Wait()
		t.SetUnActive()
	}
	if success {
		// when stop success, no more sendTimer, the t.C is empty,
		// we reset to un-active status, prevent others read from t.C
		t.SetUnActive()
	}
	return success
}

// Reset will reuse a timer with a new time start
// old timer will be stop safely
func (t *Timer) Reset(dur time.Duration) bool {
	t.Stop()
	active := t.t.Reset(dur)
	t.setActive()
	return active
}

// Go's timer has errors:
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
