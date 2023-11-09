package callwheel

import (
	"testing"
	"time"
)

func TestCallwheel(t *testing.T) {
	// naming a constructor var defaults the others?
	epoch := time.Now()
	must_continue := true
	cw := CallWheel{Size: 10}
	cw.Begin()
	cw.Insert(3, func(){ t.Log("3 sec") })
	cw.Insert(5, func() { must_continue = false })
	for must_continue {
		thisEpoch := time.Now()
		if thisEpoch.Sub(epoch).Milliseconds() >= 1000 {
			epoch = thisEpoch
			t.Log("Tick.")
			cw.Tick()
		}
	}
	cw.End()
}

