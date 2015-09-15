package main

import "os"

const PORT = "1205"

// Port returns the port used by pilad.
func Port() string {
	port := os.Getenv("PILADB_PORT")
	if port == "" {
		return PORT
	}

	return port
}
