// package pila represents the Go library that handles the Pila, databases
// and stacks.
package pila

// Pilla contains a reference to all the existing Databases, i.e.
// the currently running piladb instance
type Pila struct {
	Databases map[string]*Database
}

// NewPila return a blank piladb instance
func NewPila() *Pila {
	databases := make(map[string]*Database)
	pila := &Pila{
		Databases: databases,
	}
	return pila
}

// CreateDatabase creates a database given a name, and build the relation
// between such database and the Pila. It return the ID of the database.
func (p *Pila) CreateDatabase(name string) string {
	db := NewDatabase(name)
	db.Pila = p
	p.Databases[db.ID] = db
	return db.ID
}
