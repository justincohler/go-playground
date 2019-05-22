package stack

import (
	"fmt"
	"sync"
)

// BoundedQueue is a fixed size implementation of an IntStack.
type BoundedQueue struct {
	queue    []int
	hasSpace *sync.Cond
	hasItems *sync.Cond
	capacity int
	sync.Mutex
}

// Push adds a value to the tail of the queue.
func (q *BoundedQueue) Push(value int) {
	q.Lock()
	defer q.Unlock()
	for len(q.queue) == q.capacity {
		q.hasSpace.Wait()
	}
	q.queue = append(q.queue, value)
	if len(q.queue) > 0 {
		q.hasItems.Broadcast()
	}
	fmt.Println("Pushed", value)
}

// Pop returns (and removes) a value from the head of the queue.
func (q *BoundedQueue) Pop() int {
	q.Lock()
	defer q.Unlock()
	var res int
	for len(q.queue) == 0 {
		q.hasItems.Wait()
	}
	res = q.queue[0]
	q.queue = q.queue[1:]
	if len(q.queue) < q.capacity {
		q.hasSpace.Broadcast()
	}
	fmt.Println("Popped", res)
	return res
}

// NewBounded returns a BoundedQueue impl of IntStack.
func NewBounded(capacity int) IntStack {
	q := &BoundedQueue{capacity: capacity}
	q.hasSpace = sync.NewCond(q)
	q.hasItems = sync.NewCond(q)
	return q
}
