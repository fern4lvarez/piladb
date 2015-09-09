package pila

import (
	"fmt"

	"github.com/fern4lvarez/piladb/pkg/uuid"
)

// Database represents a piladb database
type Database struct {
	// ID is a unique identifier of the database
	ID fmt.Stringer
	// Name of the database
	Name string
	// Pointer to the current piladb instance
	Pila *Pila
}

// NewDatabase creates a new Database given a name,
// without any link to the piladb instance.
func NewDatabase(name string) *Database {
	return &Database{
		ID:   uuid.New(name),
		Name: name,
	}
}
