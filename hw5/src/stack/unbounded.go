package stack

import (
	"fmt"
	"sync"
)

// UnboundedQueue is a linked-list unbounded queue impl.
type UnboundedQueue struct {
	head    *IntNode
	tail    *IntNode
	enqLock sync.Mutex
	deqLock sync.Mutex
}

// Push adds a value to the tail of the queue.
func (q *UnboundedQueue) Push(value int) {
	q.enqLock.Lock()
	defer q.enqLock.Unlock()

	newNode := &IntNode{value: value}
	q.tail.next = newNode
	q.tail = newNode
	fmt.Println("Pushed", value)
}

// Pop returns (and removes) a value from the head of the queue.
func (q *UnboundedQueue) Pop() int {
	var res int
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
	fmt.Println("Popped", res)
	return res // will return 0 if poping empty queue
}

// NewUnbounded returns an UnboundedQueue impl of IntStack.
func NewUnbounded() IntStack {
	q := &UnboundedQueue{head: &IntNode{}}
	q.tail = q.head
	return q
}
