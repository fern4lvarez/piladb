package main

import "os"

// PORT represent the default port where pilad will running.
const PORT = "1205"

// Port returns the port used by pilad.
func Port() string {
	port := os.Getenv("PILADB_PORT")
	if port == "" {
		return PORT
	}

	return port
}
