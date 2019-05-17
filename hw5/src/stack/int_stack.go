package stack

// IntStack provides a stack for int values.
type IntStack interface {
	Push(value int)
	Pop() int
}

// IntNode provides a linked-list node style for int values.
type IntNode struct {
	value int
	next  *IntNode
}
