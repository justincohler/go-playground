package queue

import (
	"sync"
)

// Stack provides a stack for str values.
type Stack interface {
	Push(value Value)
	Pop() Value
	Empty() bool
}

// Node provides a linked-list node style for str values.
type Node struct {
	value Value
	next  *Node
}

// Value acts as a generic interface to perform queue ops on any object.
type Value interface {
}

// UnboundedQueue is a linked-list unbounded queue impl.
type UnboundedQueue struct {
	head    *Node
	tail    *Node
	enqLock sync.Mutex
	deqLock sync.Mutex
}

// Push adds a value to the tail of the queue.
func (q *UnboundedQueue) Push(value Value) {
	q.enqLock.Lock()
	defer q.enqLock.Unlock()

	newNode := &Node{value: value}
	q.tail.next = newNode
	q.tail = newNode
}

// Pop returns (and removes) a value from the head of the queue.
func (q *UnboundedQueue) Pop() Value {
	var res Value
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
func NewUnbounded() Stack {
	q := &UnboundedQueue{head: &Node{}}
	q.tail = q.head
	return q
}
