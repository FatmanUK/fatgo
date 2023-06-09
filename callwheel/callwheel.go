package callwheel

import (
	"container/ring"
	"container/list"
	"sync"
)

// Implements a simple Franta-Maly callwheel.

// The value of cws (here: size) should be kept comparable to the
// maximum possible number of outstanding timers to reduce the time
// taken to traverse the entire callwheel.
type CallWheel struct {
	Size int // aka. cws
	ring *ring.Ring
	mutex sync.Mutex // worried about sync issues
}

// TODO: possibly abstract this away to nothing? To accept any function?
type TimedOperation func()

type Event struct {
	fn TimedOperation
	ttl int
}

func (re *CallWheel) Begin() {
	re.mutex.Lock()
	defer re.mutex.Unlock()
	// create ring (size) of list
	re.ring = ring.New(re.Size)
	// initialise linked lists
	for i := re.ring.Len(); i > 0; i-- {
		re.ring.Value = list.New()
		re.ring = re.ring.Next()
	}
}

func (re *CallWheel) End() {
	//re.mutex.Lock()
	//defer re.mutex.Unlock()
}

// ticks can be seconds or minutes or whatever --- we don't care
func (re *CallWheel) Tick() {
	re.mutex.Lock()
	defer re.mutex.Unlock()
	// Move the current pointer forward one
	re.ring = re.ring.Next()
	thisList := re.ring.Value.(*list.List)
	for e := thisList.Front(); e != nil; e = e.Next() {
		thisEvent := e.Value.(*Event)
		// decrement the ttl
		//fmt.Println("B: %d", thisEvent.ttl)
		thisEvent.ttl = thisEvent.ttl - 1
		//fmt.Println("A: %d", thisEvent.ttl)
		if thisEvent.ttl == 0 {
			go thisEvent.fn()
		}
	}
	// TODO: garbage collect all ttl <= 0
}

func (re *CallWheel) Insert(num_ticks int, fn TimedOperation) {
	re.mutex.Lock()
	defer re.mutex.Unlock()
	ttl := (num_ticks / re.Size) + 1 // ttl of the timer
	index := num_ticks % re.Size // relative ring index
	//fmt.Println("Inserting at", index, "+", (re.size * ttl))
	e := &Event{fn: fn, ttl: ttl}
	for n := index; n > 0; n-- {
		re.ring = re.ring.Next()
	}
	thisList := re.ring.Value.(*list.List)
	thisList.PushBack(e)
	for n := index; n > 0; n-- {
		re.ring = re.ring.Prev()
	}
}
