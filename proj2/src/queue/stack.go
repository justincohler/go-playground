package queue

// StringStack provides a stack for str values.
type StringStack interface {
	Push(value string)
	Pop() string
	Empty() bool
}

// StringNode provides a linked-list node style for str values.
type StringNode struct {
	value string
	next  *StringNode
}
