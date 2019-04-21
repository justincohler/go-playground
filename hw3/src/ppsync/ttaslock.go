package ppsync

import (
	"sync/atomic"
)

// A TTASLock is a mutual exclusion lock that represents a
// test-and-test-and-set lock. The zero value for a TTASLock
// is an unlocked mutex.
type TTASLock struct {
	value int32
}

// Lock locks lock. If the lock is already in use, the calling goroutine
// blocks until the lock is available.
func (lock *TTASLock) Lock() {
	for {
		for lock.value == TRUE {
		}
		if atomic.CompareAndSwapInt32(&lock.value, FALSE, TRUE) {
			return
		}
	}
}

// Unlock unlocks lock.
// It is a run-time error if lock is not locked on entry to Unlock.
//
// A locked lock is not associated with a particular goroutine.
// It is allowed for one goroutine to lock a lock and then
// arrange for another goroutine to unlock it.
func (lock *TTASLock) Unlock() {
	atomic.StoreInt32(&lock.value, FALSE)
}
