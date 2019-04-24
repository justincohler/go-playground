package lock

// RWLock implements RWMutex.
type RWLock struct {
}

// Lock locks rw for writing. If the lock is already locked for reading or writing, Lock blocks until the lock is available.
func (rw *RWLock) Lock() {

}

// Unlock unlocks a Read/Write thread lock.
func (rw *RWLock) Unlock() {

}

// RLock locks rw for reading.
//
// It should not be used for recursive read locking; a blocked Lock call excludes new readers from acquiring the lock. See the documentation on the RWMutex type.
func (rw *RWLock) RLock() {

}

// RUnlock undoes a single RLock call; it does not affect other simultaneous readers.
// It is a run-time error if rw is not locked for reading on entry to RUnlock.
func (rw *RWLock) RUnlock() {

}
