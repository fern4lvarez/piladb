package stack

// Stacker represents an interface that contains all the
// required methods to implement a Stack that can be
// used in piladb.
type Stacker interface {
	// Push an element into a Stack
	Push(element interface{})
	// Pop the topmost element of a stack
	Pop() (interface{}, bool)
	// Base bases a new element at the bottom of the stack
	Base(element interface{})
	// Sweep the bottommost element of a stack
	Sweep() (interface{}, bool)
	// SweepPush sweeps the bottommost element of a stack
	// and pushes another on top
	SweepPush(element interface{}) (interface{}, bool)
	// Size returns the size of the Stack
	Size() int
	// Peek returns the topmost element of the Stack
	Peek() interface{}
	// Flush flushes a Stack
	Flush()
}
