package vars

import "fmt"

const (
	// MaxStackSize is the maximun number
	// of elements that a stack can contain.
	MaxStackSize = "MAX_STACK_SIZE"
	// MaxStackSizeDefault represents the default value
	// of MaxStackSize.
	MaxStackSizeDefault = -1

	// ReadTimeout is the maximun duration
	// before timing out the read of a request
	// to pilad.
	ReadTimeout = "READ_TIMEOUT"
	// ReadTimeoutDefault represents the default value
	// of ReadTimeout.
	ReadTimeoutDefault = 30

	// WriteTimeout is the maximun duration
	// before timing out the write of a response
	// from pilad.
	WriteTimeout = "WRITE_TIMEOUT"
	// WriteTimeoutDefault represents the default value
	// of WriteTimeout.
	WriteTimeoutDefault = 45

	// Port is the TCP port number where pilad
	// is running. Port number range is 1025-65536.
	Port = "PORT"
	// PortDefault represents the default value
	// of Port.
	PortDefault = 1205
)

// Env returns the environment variable name
// given a config name.
func Env(name string) string {
	return fmt.Sprintf("PILADB_%s", name)
}

// DefaultInt returns the default value of a config
// name of int type.
func DefaultInt(name string) int {
	switch name {
	case MaxStackSize:
		return MaxStackSizeDefault
	case ReadTimeout:
		return ReadTimeoutDefault
	case WriteTimeout:
		return WriteTimeoutDefault
	case Port:
		return PortDefault
	}
	return -1
}
