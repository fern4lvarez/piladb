package pila

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/fern4lvarez/piladb/pkg/stack"
	"github.com/fern4lvarez/piladb/pkg/uuid"
)

// Stack represents a stack entity in piladb.
type Stack struct {
	// ID is a unique identifier of the Stack
	ID fmt.Stringer

	// Name of the Stack
	Name string

	// Database associated to the Stack
	Database *Database

	// CreatedAt represents the date when the Stack was created
	CreatedAt time.Time

	// UpdatedAt represents the date when the Stack was updated for the last time.
	// This date must be updated when a Stack is created, and when receives a PUSH,
	// POP, or FLUSH operation.
	// Note that unlike CreatedAt, UpdatedAt is not triggered automatically
	// when one of these events happens, but it needs to be set by hand.
	UpdatedAt time.Time

	// ReadAt represents the date when the Stack was read for the last time.
	// This date must be updated when a Stack is created, accessed, and when it
	// receives a PUSH, POP, or FLUSH operation.
	// Note that unlike CreatedAt, ReadAt is not triggered automatically
	// when one of these events happens, but it needs to be set by hand.
	ReadAt time.Time

	// base represents the Stack data structure
	base stack.Stacker
}

// NewStack creates a new Stack given a name and a creation date,
// without an association to any Database.
func NewStack(name string, t time.Time) *Stack {
	s := &Stack{}
	s.Name = name
	s.SetID()
	s.CreatedAt = t
	s.base = stack.NewStack()
	return s
}

// Push an element on top of the Stack.
func (s *Stack) Push(element interface{}) {
	s.base.Push(element)
}

// Pop removes and returns the element on top of the Stack.
// If the Stack was empty, it returns false.
func (s *Stack) Pop() (interface{}, bool) {
	return s.base.Pop()
}

// Size returns the size of the Stack.
func (s *Stack) Size() int {
	return s.base.Size()
}

// Peek returns the element on top of the Stack.
func (s *Stack) Peek() interface{} {
	return s.base.Peek()
}

// Flush flushes the content of the Stack.
func (s *Stack) Flush() {
	s.base.Flush()
}

// Update takes a date and updates UpdateAt and ReadAt
// fields of the Stack.
func (s *Stack) Update(t time.Time) {
	s.UpdatedAt = t
	s.ReadAt = t
}

// Read takes a date and updates ReadAt field
// of the Stack.
func (s *Stack) Read(t time.Time) {
	s.ReadAt = t
}

// SetDatabase links the Stack with a given Database and
// recalculates its ID.
func (s *Stack) SetDatabase(db *Database) {
	s.Database = db
	s.SetID()
}

// SetID recalculates the id of the Stack based on its
// Database name and its own name.
func (s *Stack) SetID() {
	if s.Database != nil {
		s.ID = uuid.New(s.Database.Name + s.Name)
		return
	}

	s.ID = uuid.New(s.Name)
}

// SizeToJSON returns the size of the Stack encoded as json.
func (s *Stack) SizeToJSON() []byte {
	// Do not check error as we consider the size
	// of a stack valid for a JSON encoding.
	size, _ := json.Marshal(s.Size())
	return size
}

// Status returns the status of the Stack  in json format.
func (s *Stack) Status() StackStatus {
	status := StackStatus{}
	status.ID = s.ID.String()
	status.Name = s.Name
	status.Size = s.Size()
	status.Peek = s.Peek()
	status.CreatedAt = s.CreatedAt
	status.UpdatedAt = s.UpdatedAt
	status.ReadAt = s.ReadAt

	return status
}

// Element represents the payload of a Stack element.
type Element struct {
	Value interface{} `json:"element"`
}

// ToJSON converts an Element into JSON.
func (element Element) ToJSON() ([]byte, error) {
	return json.Marshal(element)
}

// Decode decodes json data into an Element.
func (element *Element) Decode(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(element)
}
