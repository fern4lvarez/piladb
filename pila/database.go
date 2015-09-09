package pila

// Database represents a piladb database
type Database struct {
	// ID is a unique identifier of the database
	ID string
	// Name of the database
	Name string
	// Pointer to the current piladb instance
	Pila *Pila
}

// NewDatabase creates a new Database given a name,
// without any link to the piladb instance.
func NewDatabase(name string) *Database {
	return &Database{
		ID:   name,
		Name: name,
	}
}
