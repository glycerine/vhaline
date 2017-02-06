package vhaline

import (
	"fmt"
	"time"
)

func (r *Replica) ParentMustHaveFailed() {
	f, e := r.ParentFailed()
	if !f {
		panic(e)
	}
}

func (r *Replica) ChildMustHaveFailed() {
	f, e := r.ChildFailed()
	if !f {
		panic(e)
	}
}

func (r *Replica) ParentMustNotHaveFailed() {
	f, e := r.ParentFailed()
	if f {
		panic(e)
	}
}

func (r *Replica) ChildMustNotHaveFailed() {
	f, e := r.ChildFailed()
	if f {
		panic(e)
	}
}

// ParentFailed returns true if we
// have no parent, or if we had one
// but failure was detected.
func (r *Replica) ParentFailed() (failed bool, explain string) {
	now := time.Now().UTC()
	alive, lastContact, ttl := r.parentLiveness.isAlive(now)
	explain = fmt.Sprintf("parent alive=%v, last-contact: %v, ttl: %v, now=%v. lastContact.Add(ttl).After(now)=%v.",
		alive, lastContact, ttl, now, lastContact.Add(ttl).After(now))
	failed = !alive
	return
}

// ChildFailed returns true if we
// have no child, or if we had one
// but failure was detected.
func (r *Replica) ChildFailed() (failed bool, explain string) {
	now := time.Now().UTC()
	alive, lastContact, ttl := r.childLiveness.isAlive(now)
	explain = fmt.Sprintf("child alive=%v, last-contact: %v, ttl: %v, now=%v. lastContact.Add(ttl).After(now)=%v.",
		alive, lastContact, ttl, now, lastContact.Add(ttl).After(now))
	failed = !alive
	return
}
