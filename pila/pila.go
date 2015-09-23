// Package pila represents the Go library that handles the Pila, databases
// and stacks.
package pila

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
)

// Pila contains a reference to all the existing Databases, i.e.
// the currently running piladb instance.
type Pila struct {
	Databases map[fmt.Stringer]*Database
}

// pilaStatus contains the status of the Pila instance.
type pilaStatus struct {
	NumberDatabases int      `json:"number_of_databases"`
	Databases       []string `json:"databases"`
}

// NewPila return a blank piladb instance
func NewPila() *Pila {
	databases := make(map[fmt.Stringer]*Database)
	pila := &Pila{
		Databases: databases,
	}
	return pila
}

// CreateDatabase creates a database given a name, and build the relation
// between such database and the Pila. It return the ID of the database.
func (p *Pila) CreateDatabase(name string) fmt.Stringer {
	db := NewDatabase(name)
	db.Pila = p
	p.Databases[db.ID] = db
	return db.ID
}

// AddDatabase adds a given Database to the Pila. It returns and error if the Database
// already had an assigned Pila, or if the Pila already contained the Database.
func (p *Pila) AddDatabase(db *Database) error {
	if db.Pila != nil {
		return errors.New("database already added to a pila")
	}
	if _, ok := p.Databases[db.ID]; ok {
		return errors.New("pila already contains database")
	}

	db.Pila = p
	p.Databases[db.ID] = db
	return nil
}

// RemoveDatabase deletes a Database given an ID from the Pila and returns
// true if it succeeded.
func (p *Pila) RemoveDatabase(id fmt.Stringer) bool {
	db, ok := p.Databases[id]
	if !ok {
		return false
	}

	delete(p.Databases, id)
	db.Pila = nil
	return true
}

// Database determines if a Database given by an ID is part
// of the Pila, returning a pointer to the Database and a boolean
// flag.
func (p *Pila) Database(id fmt.Stringer) (*Database, bool) {
	db, ok := p.Databases[id]
	return db, ok
}

// Status returns the status of the Pila instance in json format.
func (p *Pila) Status() []byte {
	ps := pilaStatus{}
	ps.NumberDatabases = len(p.Databases)
	dbs := make([]string, 0, len(p.Databases))
	for d := range p.Databases {
		dbs = append(dbs, d.String())
	}

	var dbsSorted sort.StringSlice = dbs
	dbsSorted.Sort()
	ps.Databases = dbsSorted

	// Do not check error as the Status type does
	// not contain types that could cause such case.
	// See http://golang.org/src/encoding/json/encode.go?s=5438:5481#L125
	b, _ := json.Marshal(ps)
	return b
}
