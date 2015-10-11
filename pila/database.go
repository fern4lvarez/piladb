package pila

import (
	"encoding/json"
	"fmt"
	"sort"

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

// databaseStatus represents the status of a Database.
type databaseStatus struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	NumberStacks int      `json:"number_of_stacks"`
	Stacks       []string `json:"stacks,omitempty"`
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
	db.Stacks[stack.Id] = stack
	return stack.Id
}

// AddStack adds a given Stack to the Database, returning
// an error if any was found.
func (db *Database) AddStack(stack *Stack) error {
	if stack.Database != nil {
		return fmt.Errorf("stack %v already added to database %v", stack.Name, stack.Database.Name)
	}

	if _, ok := db.Stacks[stack.Id]; ok {
		return fmt.Errorf("database %v already contains stack %v", db.Name, stack.Name)
	}

	stack.Database = db
	db.Stacks[stack.Id] = stack
	return nil
}

// RemoveStack removes a Stack from the Database given an id,
// returning true if it succeeded. It will return false if the
// Stack wasn't added to the Database.
func (db *Database) RemoveStack(id fmt.Stringer) bool {
	stack, ok := db.Stacks[id]
	if !ok {
		return false
	}
	stack.Database = nil
	stack.base = nil
	delete(db.Stacks, id)
	return true
}

// Status returns the status of the Database in json format.
func (db *Database) Status() []byte {
	dbs := databaseStatus{}
	dbs.ID = db.ID.String()
	dbs.Name = db.Name
	dbs.NumberStacks = len(db.Stacks)

	var ss sort.StringSlice = make([]string, len(db.Stacks))
	n := 0
	for sID := range db.Stacks {
		ss[n] = sID.String()
		n++
	}
	ss.Sort()
	dbs.Stacks = ss

	// Do not check error as the Status type does
	// not contain types that could cause such case.
	// See http://golang.org/src/encoding/json/encode.go?s=5438:5481#L125
	b, _ := json.Marshal(dbs)
	return b
}
