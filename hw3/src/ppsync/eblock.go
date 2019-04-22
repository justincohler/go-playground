package ppsync

import (
	"sync/atomic"
	"time"
)

// EBLock is a mutual exclusion lock that represents an
// Exponential Backoff lock. The zero value for an Exponential
// Backoff lock is an unlocked mutex.
type EBLock struct {
	value int32
}

// Lock locks lock. If the lock is already in use, the calling goroutine
// blocks until the lock is available.
func (lock *EBLock) Lock() {
	minDelay := 32
	maxDelay := 1024
	delay := minDelay
	for {
		for lock.value == TRUE {
		}
		if atomic.CompareAndSwapInt32(&lock.value, FALSE, TRUE) {
			return
		}
		time.Sleep(2 * time.Millisecond)
		if delay < maxDelay {
			delay *= 2
		}
	}
}

// Unlock unlocks lock.
// It is a run-time error if lock is not locked on entry to Unlock.
//
// A locked lock is not associated with a particular goroutine.
// It is allowed for one goroutine to lock a lock and then
// arrange for another goroutine to unlock it.
func (lock *EBLock) Unlock() {
	atomic.StoreInt32(&lock.value, FALSE)
}
