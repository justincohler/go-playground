package lock

import (
	"sync"
	"sync/atomic"
)

// ReadLock implements a read lock struct.
type ReadLock struct {
	mux *sync.Mutex
}

// Lock locks a given ReadLock.
func (lock *ReadLock) Lock(writer *bool, readers *int32, cond *sync.Cond) {
	lock.mux.Lock()
	defer lock.mux.Unlock()
	for *writer {
		cond.Wait()
	}
	atomic.AddInt32(readers, 1)
}

// Unlock unlocks a given ReadLock.
func (lock *ReadLock) Unlock(writer *bool, readers *int32, cond *sync.Cond) {
	lock.mux.Lock()
	defer lock.mux.Unlock()

	res := atomic.AddInt32(readers, -1)
	if res == 0 {
		cond.Broadcast()
	}
}

// WriteLock implements a write Lock struct.
type WriteLock struct {
	mux *sync.Mutex
}

// Lock locks a given WriteLock.
func (lock *WriteLock) Lock(writer *bool, readers *int32, cond *sync.Cond) {
	lock.mux.Lock()
	defer lock.mux.Unlock()
	for *readers > 0 {
		cond.Wait()
	}
	*writer = true
}

// Unlock unlocks a given WriteLock.
func (lock *WriteLock) Unlock(writer *bool, cond *sync.Cond) {
	*writer = false
	cond.Broadcast()
}

// RWLock implements RWMutex.
type RWLock struct {
	readers   int32
	writer    bool
	cond      *sync.Cond
	readLock  *ReadLock
	writeLock *WriteLock
}

// Lock locks rw for writing. If the lock is already locked for reading or writing, Lock blocks until the lock is available.
func (rw *RWLock) Lock() {
	rw.writeLock.Lock(&rw.writer, &rw.readers, rw.cond)
}

// Unlock unlocks a Read/Write thread lock.
func (rw *RWLock) Unlock() {
	rw.writeLock.Unlock(&rw.writer, rw.cond)
}

// RLock locks rw for reading.
//
// It should not be used for recursive read locking; a blocked Lock call excludes new readers from acquiring the lock. See the documentation on the RWMutex type.
func (rw *RWLock) RLock() {
	rw.readLock.Lock(&rw.writer, &rw.readers, rw.cond)
}

// RUnlock undoes a single RLock call; it does not affect other simultaneous readers.
// It is a run-time error if rw is not locked for reading on entry to RUnlock.
func (rw *RWLock) RUnlock() {
	rw.readLock.Unlock(&rw.writer, &rw.readers, rw.cond)
}
