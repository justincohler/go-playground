package main

import (
	"fmt"
	"stack"
)

func main() {

	fmt.Println("Bounded queue...")
	var bounded stack.IntStack
	bounded = stack.NewBounded(2)

	bounded.Add(7)
	go bounded.Push(1)
	go bounded.Push(2)
	go bounded.Push(3)
	go bounded.Push(4)
	go bounded.Pop()
	go bounded.Pop()
	go bounded.Pop()

	bounded.Await()

	fmt.Println("Unbounded queue...")
	var unbounded stack.IntStack
	unbounded = stack.NewUnbounded()

	unbounded.Add(7)
	go unbounded.Push(1)
	go unbounded.Push(2)
	go unbounded.Push(3)
	go unbounded.Push(4)
	go unbounded.Pop()
	go unbounded.Pop()
	go unbounded.Pop()

	unbounded.Await()
}
