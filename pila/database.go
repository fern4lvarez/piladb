package pila

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/fern4lvarez/piladb/pkg/stack"
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

// CreateStack creates a new Stack, given a name and a creation date,
// which is associated to the Database.
func (db *Database) CreateStack(name string, t time.Time) fmt.Stringer {
	return db.CreateStackWithBase(name, t, stack.NewStack())
}

// CreateStackWithBase creates a new Stack, given a name, a creation date,
// and a stack.Stacker base implementation, which is associated to the Database.
func (db *Database) CreateStackWithBase(name string, t time.Time, base stack.Stacker) fmt.Stringer {
	stack := NewStackWithBase(name, t, base)
	stack.SetDatabase(db)
	db.Stacks[stack.ID] = stack
	return stack.ID
}

// AddStack adds a given Stack to the Database, returning
// an error if any was found.
func (db *Database) AddStack(stack *Stack) error {
	if stack.Database != nil {
		return fmt.Errorf("stack %v already added to database %v", stack.Name, stack.Database.Name)
	}

	stack.SetDatabase(db)
	if _, ok := db.Stacks[stack.ID]; ok {
		stack.Database = nil
		return fmt.Errorf("database %v already contains stack %v", db.Name, stack.Name)
	}

	db.Stacks[stack.ID] = stack
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

// Status returns the status of the Database.
func (db *Database) Status() DatabaseStatus {
	dbs := DatabaseStatus{}
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

	return dbs
}

// StacksStatus returns the status of the Stacks of Database.
func (db *Database) StacksStatus() StacksStatus {
	var n int
	ss := make([]StackStatus, len(db.Stacks))
	for _, s := range db.Stacks {
		s.CreatedAt = s.CreatedAt.Local()
		s.UpdatedAt = s.UpdatedAt.Local()
		s.ReadAt = s.ReadAt.Local()
		ss[n] = s.Status()
		n++
	}

	status := StacksStatus{Stacks: ss}
	sort.Sort(status)
	return status
}

// StacksKV returns the status of the Stacks of Database
// in a key-value format.
func (db *Database) StacksKV() StacksKV {
	kv := make(map[string]interface{})
	for _, s := range db.Stacks {
		kv[s.Name] = s.Peek()
	}

	stacksKV := StacksKV{Stacks: kv}
	return stacksKV
}

// DatabaseStatus represents the status of a Database.
type DatabaseStatus struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	NumberStacks int      `json:"number_of_stacks"`
	Stacks       []string `json:"stacks,omitempty"`
}

// ToJSON converts a DatabaseStatus into JSON.
func (databaseStatus DatabaseStatus) ToJSON() []byte {
	// Do not check error as the Status type does
	// not contain types that could cause such case.
	// See http://golang.org/src/encoding/json/encode.go?s=5438:5481#L125
	b, _ := json.Marshal(databaseStatus)
	return b
}
