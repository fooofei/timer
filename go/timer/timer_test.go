package timer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
	"time"
)

// "github.com/shoenig/test" 还不够好用，失败后，不打印哪个是 expect 只看到不相登

func reuseTimerWithLength(t *time.Timer, d time.Duration) {
	if !t.Stop() && len(t.C) > 0 {
		<-t.C
	}
	t.Reset(d)
}

func reuseTimerWithDefault(t *time.Timer, d time.Duration) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
	t.Reset(d)
}

func waitTimerWithDefault() error {
	tmr := time.NewTimer(0)
	reuseTimerWithDefault(tmr, time.Minute)
	select {
	case <-tmr.C:
		return fmt.Errorf("unexpected firing of Timer")
	default:
		return nil
	}
}

func waitTimerWithLength() error {
	tmr := time.NewTimer(0)
	reuseTimerWithLength(tmr, time.Minute)
	select {
	case <-tmr.C:
		return fmt.Errorf("unexpected firing of Timer")
	default:
		return nil
	}
}

// TestErrorUseWaitTimerWithDefault 证明 select default 也不能正确使用 timer.Reset
func TestErrorUseWaitTimerWithDefault(t *testing.T) {
	runtime.GOMAXPROCS(2)
	for i := 0; ; i++ {
		if err := waitTimerWithDefault(); err != nil {
			t.Logf("occured error '%v' when loop count %v", err, i)
			break
		}
	}

}

// TestErrorUseWaitTimerWithLength 证明 len(timer.C) 也不能正确使用 timer.Reset
func TestErrorUseWaitTimerWithLength(t *testing.T) {
	runtime.GOMAXPROCS(2)
	for i := 0; ; i++ {
		if err := waitTimerWithLength(); err != nil {
			t.Logf("occured error '%v' when loop count %v", err, i)
			break
		}
	}
}

func TestMyTimerReset(tst *testing.T) {
	var t = New(time.Second)
	defer t.Stop()

	tst.Logf("wait 1s")
	const myTimerArrive = "my timer arrived"
	const myTimerNotArrive = "my timer not arrived"
	var whichArrive string
	select {
	case <-t.After(2 * time.Second):
		t.SetUnActive()
		whichArrive = myTimerArrive
	case <-time.After(time.Second):
		whichArrive = myTimerNotArrive
		// when this done, our timer is not triggered, and will set un-active for next use.
	}
	assert.Equal(tst, myTimerNotArrive, whichArrive)

	tst.Logf("wait 2s")
	whichArrive = ""

	select {
	case <-t.After(2 * time.Second):
		t.SetUnActive()
		whichArrive = myTimerArrive
	case <-time.After(4 * time.Second):
		whichArrive = myTimerNotArrive
	}
	assert.Equal(tst, myTimerArrive, whichArrive)

	tst.Logf("wait 1s")
	whichArrive = ""
	select {
	case <-t.After(3 * time.Second):
		t.SetUnActive()
		whichArrive = myTimerArrive
	case <-time.After(time.Second):
		whichArrive = myTimerNotArrive
	}
	assert.Equal(tst, myTimerNotArrive, whichArrive)
}

func TestMyTimerStopReuse(tst *testing.T) {
	var t = New(time.Second)
	defer t.Stop()

	t.Stop()
	const myTimerArrive = "my timer arrived"
	const myTimerNotArrive = "my timer not arrived"
	var whichArrive string
	select {
	case <-t.After(time.Second):
		t.SetUnActive()
		whichArrive = myTimerArrive
	case <-time.After(5 * time.Second):
		whichArrive = myTimerNotArrive
	}
	assert.Equal(tst, myTimerArrive, whichArrive)

}
