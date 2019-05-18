package main

import (
	"fmt"
	"stack"
	"sync"
)

func push(wg *sync.WaitGroup, q stack.IntStack, value int) {
	defer wg.Done()
	q.Push(value)
}

func pop(wg *sync.WaitGroup, q stack.IntStack) int {
	defer wg.Done()
	return q.Pop()
}

func main() {

	var wg sync.WaitGroup
	fmt.Println("Bounded queue...")
	var bounded stack.IntStack
	bounded = stack.NewBounded(2)

	wg.Add(8)
	go push(&wg, bounded, 1)
	go push(&wg, bounded, 2)
	go push(&wg, bounded, 3)
	go push(&wg, bounded, 4)
	go pop(&wg, bounded)
	go pop(&wg, bounded)
	go pop(&wg, bounded)
	go pop(&wg, bounded)
	wg.Wait()

	fmt.Println("Unbounded queue...")
	var unbounded stack.IntStack
	unbounded = stack.NewUnbounded()

	wg.Add(8)
	go push(&wg, unbounded, 1)
	go push(&wg, unbounded, 2)
	go push(&wg, unbounded, 3)
	go push(&wg, unbounded, 4)
	go pop(&wg, unbounded)
	go pop(&wg, unbounded)
	go pop(&wg, unbounded)
	go pop(&wg, unbounded)

	wg.Wait()
}
