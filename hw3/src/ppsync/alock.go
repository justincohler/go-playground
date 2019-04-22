package ppsync

import (
	"sync/atomic"
)

// ALock is a mutual exclusion lock that represents an
// Anderson Queue lock. The zero value for a Anderson lock
// is an unlocked mutex.
type ALock struct {
	Flags []bool
	slot  int32
	next  int32
	N     int32
}

// Lock locks lock. If the lock is already in use, the calling goroutine
// blocks until the lock is available.
func (lock *ALock) Lock() {
	slot := atomic.AddInt32(&lock.next, 1) % lock.N
	for !lock.Flags[slot%lock.N] {
	}
	lock.Flags[slot] = false
}

// Unlock unlocks lock.
// It is a run-time error if lock is not locked on entry to Unlock.
//
// A locked lock is not associated with a particular goroutine.
// It is allowed for one goroutine to lock a lock and then
// arrange for another goroutine to unlock it.
func (lock *ALock) Unlock() {
	lock.Flags[(lock.next+1)%lock.N] = true
}

// NewALock initializes and returns a new ALock
// for a given number of threads.
func NewALock(nThreads int) *ALock {
	lock := &ALock{Flags: make([]bool, nThreads), N: int32(nThreads)}
	lock.Flags[0] = true

	return lock
}
