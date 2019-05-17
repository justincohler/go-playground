package main

import (
	"fmt"
	"stack"
)

func main() {
	var q stack.IntStack
	q = stack.NewUnbounded()

	q.Push(3)
	q.Push(2)
	fmt.Println(q.Pop())
	fmt.Println(q.Pop())
	fmt.Println(q.Pop())
}
