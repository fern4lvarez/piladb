package vars

import "fmt"

const (
	// MaxStackSize is the maximun number
	// of elements that a stack can contain.
	MaxStackSize = "MAX_STACK_SIZE"

	// ReadTimeout is the maximun duration
	// before timing out the read of a request
	// to pilad.
	ReadTimeout = "READ_TIMEOUT"

	// WriteTimeout is the maximun duration
	// before timing out the write of a response
	// from pilad.
	WriteTimeout = "WRITE_TIMEOUT"
)

// Env returns the environment variable name
// given a config name.
func Env(name string) string {
	return fmt.Sprintf("PILADB_%s", name)
}
