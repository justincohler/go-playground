package stack

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

// UnboundedLockFreeQueue is a linked-list unbounded queue impl.
type UnboundedLockFreeQueue struct {
	head *IntNode
	tail *IntNode
}

// Push adds a value to the tail of the queue.
func (q *UnboundedLockFreeQueue) Push(value int) {
	newNode := &IntNode{value: value}
	for {
		last := *q.tail
		next := *last.next

		if last == *q.tail {
			if &next == nil {
				if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&last)),
					unsafe.Pointer(&next),
					unsafe.Pointer(&newNode)) {
					atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(q.tail)),
						unsafe.Pointer(newNode))
					return
				}
			} else {
				atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(q.tail)),
					unsafe.Pointer(newNode))
			}
		}
	}
	fmt.Println("Pushed", value)
}

// Pop returns (and removes) a value from the head of the queue.
func (q *UnboundedLockFreeQueue) Pop() int {
	for {
		first := *q.head
		last := *q.tail
		next := *first.next
		if first == *q.head {
			if first == last {
				if &next == nil {
					panic("Empty queue")
				}
				atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&last)),
					unsafe.Pointer(&next))
			} else {
				value := next.value
				if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(q.head)),
					unsafe.Pointer(&first),
					unsafe.Pointer(&next)) {
					fmt.Println("Popped", value)
					return value
				}
			}
		}
	}
}

// NewUnboundedLockFreeQueue returns an UnboundedLockFreeQueue impl of IntStack.
func NewUnboundedLockFreeQueue() IntStack {
	q := &UnboundedQueue{head: &IntNode{}}
	q.tail = q.head
	return q
}
