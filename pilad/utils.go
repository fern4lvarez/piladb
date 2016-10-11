package main

import (
	"fmt"
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

// MemOutput returns a formatted string given an amount
// of memory as an integer. It will print the result in B,
// KiB, MiB or GiB.
// 1KiB = 1,024 Bytes
// 1MiB = 1,048,576 Bytes
// 1GiB = 1,073,741,824 Bytes
func MemOutput(mem uint64) string {
	var memF64 = float64(mem)

	switch {
	case mem < 1024:
		return fmt.Sprintf("%dB", mem)
	case mem < 1048576:
		return fmt.Sprintf("%.2fKiB", memF64/1024)
	case mem < 1073741824:
		return fmt.Sprintf("%.2fMiB", memF64/1048576)
	default:
		return fmt.Sprintf("%.2fGiB", memF64/1073741824)
	}
}
