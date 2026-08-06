package callwheel

import (
	"fmt"
	"testing"
	"time"
)

func TestCallwheel(t *testing.T) {
	epoch := time.Now()
	must_continue := true
	cw := CallWheelFactory(10)
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
}

var DhcpOfferTimeouts = CallWheelFactory(10)

func TestInsertHang(t *testing.T) {
	mustContinue := true
	epoch := time.Now()
	{
		s := "abc123"
		fmt.Println("start testInsertHang")
		DhcpOfferTimeouts.Insert(3, func(){
			fmt.Println("Rescinded offer to", s)
		})
		DhcpOfferTimeouts.Insert(5, func() { mustContinue = false })
	}
	fmt.Println("end testInsertHang")
	for mustContinue {
		thisEpoch := time.Now()
		if thisEpoch.Sub(epoch).Milliseconds() >= 1000 {
			epoch = thisEpoch
			t.Log("Tick.")
			DhcpOfferTimeouts.Tick()
		}
	}
}
