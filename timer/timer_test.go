package timer

import (
	"testing"
	"time"
)

func TestLeftTime(t *testing.T) {
	initialTimeout := 10
	timeNanoseconds := time.Second * time.Duration(initialTimeout)
	tim := InitHourglass(timeNanoseconds)

	<-time.After(time.Second * time.Duration(1))
	if tim.LeftNanoseconds() < time.Second*time.Duration(8) ||
		tim.LeftNanoseconds() > time.Second*time.Duration(9) {
		t.Errorf("Testing timer.GetLeftTime for 1 sec err, TimeLeft: %d", int64(tim.LeftNanoseconds()))
	}

	<-time.After(time.Second * time.Duration(3))
	if tim.LeftNanoseconds() < time.Second*time.Duration(5) ||
		tim.LeftNanoseconds() > time.Second*time.Duration(6) {
		t.Errorf("Testing timer.GetLeftTime for 3 sec err, TimeLeft: %d", int64(tim.LeftNanoseconds()))
	}
}

func TestOverflow(t *testing.T) {
	initialTimeout := 10
	tim := InitHourglass(time.Second * time.Duration(initialTimeout))

	<-time.After(time.Second * time.Duration(11))
	if int(tim.LeftNanoseconds()) > 0 {
		t.Errorf("Testing timer.GetLeftTime overflow err, TimeLeft: %d", int64(tim.LeftNanoseconds()))
	}
}
