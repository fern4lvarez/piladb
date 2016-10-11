// Package stack provides a basic implementation
// of a stack using a linked list.
package stack

import "sync"

// Stack represents the stack data structure as a linked list,
// containing a pointer to the first Frame as a head and the
// size of the stack. It also contain a mutex to lock and unlock
// the access to the stack on I/O.
type Stack struct {
	head *frame
	size int
	mux  sync.Mutex
}

// frame represents an element of the stack. It contains
// data and the link to the next Frame as a pointer.
type frame struct {
	data interface{}
	next *frame
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
	s.mux.Lock()
	defer s.mux.Unlock()

	head := &frame{
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
	s.mux.Lock()
	defer s.mux.Unlock()

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

// Flush flushes the content of the stack.
func (s *Stack) Flush() {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.size = 0
	s.head = nil
}
