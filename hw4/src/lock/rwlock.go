package lock

import (
	"sync"
)

// RWLock implements RWMutex.
type RWLock struct {
	readers int32
	writer  bool
	mux     sync.Mutex
}

// Lock locks rw for writing. If the lock is already locked for reading or writing, Lock waits until readers have finished.
func (rw *RWLock) Lock() {
	for {
		rw.mux.Lock()
		if rw.readers > 0 {
			rw.mux.Unlock() // others still reading, try again
		} else {
			return // no more readers, keep lock alive until Unlock
		}
	}
}

// Unlock unlocks a Read/Write thread lock.
func (rw *RWLock) Unlock() {
	rw.mux.Unlock()
}

// RLock locks rw for reading.
//
// It should not be used for recursive read locking; a blocked Lock call excludes new readers from acquiring the lock. See the documentation on the RWMutex type.
func (rw *RWLock) RLock() {
	rw.mux.Lock()
	defer rw.mux.Unlock()
	rw.readers++
}

// RUnlock undoes a single RLock call; it does not affect other simultaneous readers.
// It is a run-time error if rw is not locked for reading on entry to RUnlock.
func (rw *RWLock) RUnlock() {
	rw.mux.Lock()
	defer rw.mux.Unlock()
	rw.readers--
}
