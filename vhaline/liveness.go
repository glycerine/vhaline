package vhaline

import (
	"fmt"
	"sync"
	"time"

	"github.com/glycerine/idem"
)

type liveness struct {
	lastContact time.Time
	ttl         time.Duration
	mut         sync.Mutex
	alive       bool
}

func newLiveness(ttl time.Duration) liveness {
	return liveness{ttl: ttl}
}

func (s *liveness) isAlive(now time.Time) (alive bool, lastContact time.Time, ttl time.Duration) {
	s.mut.Lock()
	ttl = s.ttl
	lastContact = s.lastContact
	alive = lastContact.Add(ttl).After(now)
	s.alive = alive
	s.mut.Unlock()
	return
}

func (s *liveness) heardFrom(now time.Time) {
	s.mut.Lock()
	if s.lastContact.Before(now) {
		s.lastContact = now.UTC()
	}
	s.alive = true
	s.mut.Unlock()
}

func (s *liveness) reset() {
	s.mut.Lock()
	s.lastContact = time.Time{}
	s.alive = false
	s.mut.Unlock()
}

type waituntil struct {
	ready     *idem.IdemCloseChan
	readylist []*idem.IdemCloseChan
	matchme   int64
}

func (w *waituntil) String() string {
	return fmt.Sprintf(" waituntil{matchme: %v} \n", w.matchme)
}

func newWaitUntil() *waituntil {
	return &waituntil{
		ready: idem.NewIdemCloseChan(),
	}
}

type healthCheckCounter struct {
	mut            sync.Mutex
	next           int64
	waitlist       *wtree
	systemShutdown chan bool // avoid deadlock on shutdown
	ttl            time.Duration
	replica        *Replica // up to containing struct
}

func newHealthCheckCounter(
	systemStopping chan bool, ttl time.Duration, replica *Replica,
) *healthCheckCounter {
	r := &healthCheckCounter{
		systemShutdown: systemStopping,
		waitlist:       newWtree(),
		ttl:            ttl,
		replica:        replica,
	}
	return r
}

// internalWakeWithoutLock is for internal
// use only, it is not not goroutine safe; caller
// must be holding the lock on s.
// Wake everyone who is waiting strictly
// before s.next.
func (s *healthCheckCounter) internalWakeWithoutLock() {
	if s.waitlist.Len() == 0 {
		return
	}
	s.waitlist.deleteThrough(s.next-1, func(goner *waituntil, through int64) {
		goner.ready.Close()
		for _, g := range goner.readylist {
			g.Close()
		}
	})
}

func (s *healthCheckCounter) lastNum() int64 {
	s.mut.Lock()
	n := s.next - 1
	s.mut.Unlock()
	return n
}

// waitUntil blocks until the s.next is incremented
// to matchme or greater by the regular heartbeats.
//
// this is inherently a little racy. but we
// use it only in tests. Should never be used
// in non-test code.
func (s *healthCheckCounter) waitUntil(matchme int64) error {
	s.mut.Lock()
	wu := newWaitUntil()
	wu.matchme = matchme
	already := s.waitlist.get(wu.matchme)
	if already == nil {
		s.waitlist.insert(wu)
	} else {
		already.readylist = append(already.readylist, wu.ready)
	}
	s.mut.Unlock()
	select {
	case <-s.systemShutdown:
		return fmt.Errorf("shutting down")
	case <-wu.ready.Chan:
		//p("5555555 ready with wu.matchme=%v\n", wu.matchme)
		return nil
	}
}

// beat advances and returns the current heartbeat count.
// It wakes up any goroutines waiting until a particular
// heartbeat to happen.
func (s *healthCheckCounter) beat() int64 {
	s.mut.Lock()
	n := s.next
	s.next++
	s.internalWakeWithoutLock()
	s.mut.Unlock()
	return n
}
