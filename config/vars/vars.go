package vars

import "fmt"

const (
	// MaxStackSize is the maximun number
	// of elements that a stack can contain.
	MaxStackSize = "MAX_STACK_SIZE"
)

// Env returns the environment variable name
// given a config name.
func Env(name string) string {
	return fmt.Sprintf("PILADB_%s", name)
}
