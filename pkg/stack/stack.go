// Package stack provides a basic implementation
// of a stack using a linked list.
package stack

import "sync"

// Stack implements the Stacker interface, and represents the stack
// data structure as a linked list, containing a pointer
// to the first Frame as a head, the last one as a base,
// and the size of the stack.
// It also contain a mutex to lock and unlock
// the access to the stack at I/O operations.
type Stack struct {
	head *frame
	tail *frame
	size int
	mux  sync.RWMutex
}

// frame represents an element of the stack. It contains
// data, and links to the up and down frames as pointers.
type frame struct {
	data interface{}
	down *frame
	up   *frame
}

// NewStack returns a blank stack, where head is nil and size
// is 0.
func NewStack() *Stack {
	return &Stack{}
}

// Push adds a new element on top of the stack, creating
// a new head holding this data and updating its head on
// top of the stack's head. It will update the tail
// only if the stack was empty.
func (s *Stack) Push(element interface{}) {
	s.mux.Lock()
	defer s.mux.Unlock()

	head := &frame{
		data: element,
		down: s.head,
		up:   nil,
	}

	// connect former head with new head
	// if Stack had already an element
	if s.size > 0 {
		s.head.up = head
	}

	s.head = head
	s.size++

	// tail and head are the same element
	// when pushing a first one
	if s.size == 1 {
		s.tail = head
		return
	}

	// update the tail when the pushed
	// element is the only on top of the
	// tail
	if s.size == 2 {
		s.tail.up = s.head
	}
}

// Pop removes and returns the element on top of the stack,
// updating its head to the Frame underneath. If the stack was empty,
// it returns false.
func (s *Stack) Pop() (interface{}, bool) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.size == 0 {
		return nil, false
	}

	element := s.head.data
	s.head = s.head.down
	s.size--

	// update the tail when it's the
	// only element after the Pop operation
	if s.size == 1 {
		s.tail.up = nil
	}

	return element, true
}

// Base bases the stack on top of a new element, so
// this element becomes the bottommost element of the
// stack.
func (s *Stack) Base(element interface{}) {
	s.mux.Lock()
	defer s.mux.Unlock()

	tail := &frame{
		data: element,
		down: nil,
		up:   s.tail,
	}

	s.tail = tail
	s.size++

	// tail and head are the same element
	// when basing a first one
	if s.size == 1 {
		s.head = tail
		return
	}

	// link the element above with new
	// tail
	if s.size > 1 {
		s.tail.up.down = tail
	}
}

// Sweep removes and returns the element at the bottom of the stack,
// turning the Frame above into the new tail. If the stack was empty,
// it returns false.
func (s *Stack) Sweep() (interface{}, bool) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.size == 0 {
		return nil, false
	}

	element := s.tail.data
	s.size--

	// head and tail are nil
	// if Stack has no elements
	if s.size == 0 {
		s.head = nil
		s.tail = nil
		return element, true
	}

	// head becomes the tail when
	// is the remaining element in Stack
	if s.size == 1 {
		s.head.down = nil
		s.head.up = nil
		s.tail = s.head
		return element, true
	}

	s.tail = s.tail.up
	s.tail.down = nil
	return element, true
}

// SweepPush makes a Sweep and Push operations as an atomic
// one, returning the swept element and the validity of
// the operation. If Stack was empty, it will return false.
func (s *Stack) SweepPush(element interface{}) (interface{}, bool) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.size == 0 {
		return nil, false
	}

	head := &frame{
		data: element,
		down: s.head,
		up:   nil,
	}

	swept := s.tail.data

	// head becomes tail when there is
	// a single element in the Stack
	if s.size == 1 {
		head.down = nil
		s.tail = head
		s.head = head
		return swept, true
	}

	// set head
	s.head.up = head
	s.head = head

	// sweep tail
	s.tail = s.tail.up
	s.tail.down = nil

	return swept, true
}

// Size returns the number of elements that a stack contains.
func (s *Stack) Size() int {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.size
}

// Peek returns the element on top of the stack.
func (s *Stack) Peek() interface{} {
	s.mux.RLock()
	defer s.mux.RUnlock()

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
	s.tail = nil
}
