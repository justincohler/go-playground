package stack

import (
	"fmt"
	"sync"
)

// UnboundedQueue is a linked-list unbounded queue impl.
type UnboundedQueue struct {
	wg      sync.WaitGroup
	head    *IntNode
	tail    *IntNode
	enqLock sync.Mutex
	deqLock sync.Mutex
}

// Push adds a value to the tail of the queue.
func (q *UnboundedQueue) Push(value int) {
	defer q.wg.Done()
	q.enqLock.Lock()
	defer q.enqLock.Unlock()

	newNode := &IntNode{value: value}
	q.tail.next = newNode
	q.tail = newNode
	fmt.Println("Pushed", value)
}

// Pop returns (and removes) a value from the head of the queue.
func (q *UnboundedQueue) Pop() int {
	defer q.wg.Done()
	var res int
	q.deqLock.Lock()
	defer q.deqLock.Unlock()

	if q.head.next != nil {
		res = q.head.next.value
		q.head = q.head.next
	}
	fmt.Println("Popped", res)
	return res // will return 0 if poping empty queue
}

// Await allows queues to be completed.
func (q *UnboundedQueue) Await() {
	q.wg.Wait()
}

// Add allows routines added to the waitGroup.
func (q *UnboundedQueue) Add(routines int) {
	q.wg.Add(routines)
}

// NewUnbounded returns an UnboundedQueue impl of IntStack.
func NewUnbounded() IntStack {
	q := &UnboundedQueue{head: &IntNode{}}
	q.tail = q.head
	return q
}
