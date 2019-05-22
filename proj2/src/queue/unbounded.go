package queue

import (
	"sync"
)

// UnboundedQueue is a linked-list unbounded queue impl.
type UnboundedQueue struct {
	head    *StringNode
	tail    *StringNode
	enqLock sync.Mutex
	deqLock sync.Mutex
}

// Push adds a value to the tail of the queue.
func (q *UnboundedQueue) Push(value string) {
	q.enqLock.Lock()
	defer q.enqLock.Unlock()

	newNode := &StringNode{value: value}
	q.tail.next = newNode
	q.tail = newNode
}

// Pop returns (and removes) a value from the head of the queue.
func (q *UnboundedQueue) Pop() string {
	var res string
	for {
		q.deqLock.Lock()
		if q.head.next == nil {
			q.deqLock.Unlock()
		} else {
			break
		}
	}
	defer q.deqLock.Unlock()

	res = q.head.next.value
	q.head = q.head.next
	return res // will return '' if poping empty queue
}

// Empty returns true if the queue is empty
func (q *UnboundedQueue) Empty() bool {
	q.enqLock.Lock()
	q.deqLock.Lock()
	defer q.enqLock.Unlock()
	defer q.deqLock.Unlock()
	if q.head.next == nil {
		return true
	}
	return false
}

// NewUnbounded returns an UnboundedQueue impl of IntStack.
func NewUnbounded() StringStack {
	q := &UnboundedQueue{head: &StringNode{}}
	q.tail = q.head
	return q
}
