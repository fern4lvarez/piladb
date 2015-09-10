package pila

import (
	"fmt"

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
	// Base represent the Stack data structure
	Base *stack.Stack
}

// NewStack creates a new Stack given a name without
// an association to any Database.
func NewStack(name string) *Stack {
	s := &Stack{}
	s.ID = uuid.New(name)
	s.Name = name
	s.Base = stack.NewStack()
	return s
}
