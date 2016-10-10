package main

import (
	"runtime"

	"github.com/fern4lvarez/piladb/pila"
	"github.com/fern4lvarez/piladb/pkg/uuid"
)

// ResourceDatabase will return the right Database resource
// given a Conn and a database ID or Name.
func ResourceDatabase(conn *Conn, databaseInput string) (*pila.Database, bool) {
	db, ok := conn.Pila.Database(uuid.UUID(databaseInput))
	if !ok {
		// Fallback to find by database name
		db, ok = conn.Pila.Database(uuid.New(databaseInput))
	}

	return db, ok
}

// ResourceStack will return the right Stack resource
// given a Database and a Stack ID or Name.
func ResourceStack(db *pila.Database, stackInput string) (*pila.Stack, bool) {
	stack, ok := db.Stacks[uuid.UUID(stackInput)]
	if !ok {
		// Fallback to find by stack name
		stack, ok = db.Stacks[uuid.New(db.Name+stackInput)]
	}

	return stack, ok
}

// MemStats fetches the memory statistics provided
// by the Go stdlib.
func MemStats() *runtime.MemStats {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return &mem
}
