// Package stack provides a basic implementation
// of a stack using a linked list.
package stack

// Frame represents an element of the stack. It contains
// data and the link to the next Frame as a pointer.
type Frame struct {
	data interface{}
	next *Frame
}

// Stack represents the stack data structure as a linked list,
// containing a pointer to the first Frame as a head and the
// size of the stack.
type Stack struct {
	head *Frame
	size int
}

// NewStack returns a blank stack, where head is nil and size
// is 0.
func NewStack() *Stack {
	return &Stack{}
}

// Push adds a new element on top of the stack, creating
// a new head holding this data and updating its head to
// the previous stack's head.
func (s *Stack) Push(element interface{}) {
	head := &Frame{
		data: element,
		next: s.head,
	}
	s.head = head
	s.size++
}

// Pop removes and returns the element on top of the stack,
// updating its head to the next Frame. If the stack was empty,
// it returns false.
func (s *Stack) Pop() (interface{}, bool) {
	if s.head == nil {
		return nil, false
	}

	element := s.head.data
	s.head = s.head.next
	s.size--
	return element, true
}

// Size returns the number of elements that a stack contains.
func (s *Stack) Size() int {
	return s.size
}

// Peek returns the element on top of the stack.
func (s *Stack) Peek() interface{} {
	if s.head == nil {
		return nil
	}
	return s.head.data
}
