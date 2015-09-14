package pila

import (
	"fmt"

	"github.com/fern4lvarez/piladb/pkg/uuid"
)

// Database represents a piladb database.
type Database struct {
	// ID is a unique identifier of the database
	ID fmt.Stringer
	// Name of the database
	Name string
	// Pointer to the current piladb instance
	Pila *Pila
	// Stacks associated to Database mapped by their ID
	Stacks map[fmt.Stringer]*Stack
}

// NewDatabase creates a new Database given a name,
// without any link to the piladb instance.
func NewDatabase(name string) *Database {
	stacks := make(map[fmt.Stringer]*Stack)
	return &Database{
		ID:     uuid.New(name),
		Name:   name,
		Stacks: stacks,
	}
}

// CreateStack creates a new Stack, given a name, which is associated
// to the Database.
func (db *Database) CreateStack(name string) fmt.Stringer {
	stack := NewStack(name)
	stack.Database = db
	db.Stacks[stack.ID] = stack
	return stack.ID
}
