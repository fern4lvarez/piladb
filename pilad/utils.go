package main

import (
	"github.com/fern4lvarez/piladb/pila"
	"github.com/fern4lvarez/piladb/pkg/uuid"
)

// ResourceDatabase will return the right Database resource]
// given a Conn and a database ID or Name.
func ResourceDatabase(conn *Conn, databaseInput string) (*pila.Database, bool) {
	db, ok := conn.Pila.Database(uuid.UUID(databaseInput))
	if !ok {
		// Fallback to find by database name
		db, ok = conn.Pila.Database(uuid.New(databaseInput))
	}

	return db, ok
}
