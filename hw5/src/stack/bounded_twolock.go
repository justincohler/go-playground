package stack

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// BoundedQueue is a fixed size implementation of an IntStack.
type BoundedTwoLockQueue struct {
	queue    []int
	hasSpace *sync.Cond
	hasItems *sync.Cond
	capacity int
	size     int32
	headLock sync.Mutex
	tailLock sync.Mutex
}

// Push adds a value to the tail of the queue.
func (q *BoundedTwoLockQueue) Push(value int) {
	q.tailLock.Lock()
	for len(q.queue) == q.capacity {
		q.hasSpace.Wait()
	}
	size := atomic.AddInt32(&q.size, 1)
	q.queue = append(q.queue, value)
	q.tailLock.Unlock()
	if size > 0 {
		q.headLock.Lock()
		q.hasItems.Broadcast()
		q.headLock.Unlock()
	}
	fmt.Println("Pushed", value)
}

// Pop returns (and removes) a value from the head of the queue.
func (q *BoundedTwoLockQueue) Pop() int {
	q.headLock.Lock()
	var res int
	for q.size == 0 {
		q.hasItems.Wait()
	}
	res = q.queue[0]
	q.queue = q.queue[1:]
	size := atomic.AddInt32(&q.size, -1)
	q.headLock.Unlock()
	if size < int32(q.capacity) {
		q.tailLock.Lock()
		q.hasSpace.Broadcast()
		q.tailLock.Unlock()
	}
	fmt.Println("Popped", res)
	return res
}

// NewBoundedTwoLockQueue returns a BoundedTwoLockQueue impl of IntStack.
func NewBoundedTwoLockQueue(capacity int) IntStack {
	q := &BoundedTwoLockQueue{capacity: capacity}
	q.hasSpace = sync.NewCond(&q.tailLock)
	q.hasItems = sync.NewCond(&q.headLock)
	return q
}
