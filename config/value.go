package config

import (
	"strconv"
	"time"

	"github.com/fern4lvarez/piladb/config/vars"
)

// MaxStackSize returns the value of MAX_STACK_SIZE.
// Type: int, Default: -1
func (c *Config) MaxStackSize() int {
	maxSize := c.Get(vars.MaxStackSize)
	return intValue(maxSize, vars.MaxStackSizeDefault)
}

// ReadTimeout returns the value of READ_TIMEOUT.
// Type: time.Duration, Default: 30
func (c *Config) ReadTimeout() time.Duration {
	readTimeout := c.Get(vars.ReadTimeout)
	t := intValue(readTimeout, vars.ReadTimeoutDefault)
	return time.Duration(t)
}

// WriteTimeout returns the value of WRITE_TIMEOUT.
// Type: time.Duration, Default: 45
func (c *Config) WriteTimeout() time.Duration {
	writeTimeout := c.Get(vars.WriteTimeout)
	t := intValue(writeTimeout, vars.WriteTimeoutDefault)
	return time.Duration(t)
}

// Port returns the value of PORT.
// Type: int, Default: 1205
func (c *Config) Port() int {
	port := c.Get(vars.Port)
	t := intValue(port, vars.PortDefault)

	if t < 1025 || t > 65536 {
		return vars.PortDefault
	}
	return t
}

// PushWhenFull returns the value of PUSH_WHEN_FULL.
// Type: bool, Default: false
func (c *Config) PushWhenFull() bool {
	pushWhenFull := c.Get(vars.PushWhenFull)
	return boolValue(pushWhenFull, vars.PushWhenFullDefault)
}

// intValue returns an Integer value given another value as an
// interface. If conversion fails, a default value is used.
func intValue(value interface{}, defaultValue int) int {
	switch value.(type) {
	case int:
		if i := value.(int); i < 0 {
			return defaultValue
		}
		return value.(int)
	case float64:
		return int(value.(float64))
	case string:
		i, err := strconv.Atoi(value.(string))
		if err != nil {
			return defaultValue
		}
		return i
	default:
		return defaultValue
	}
}

// boolValue returns a Boolean value given another value as an
// interface. If conversion fails, a default value is used.
func boolValue(value interface{}, defaultValue bool) bool {
	switch value.(type) {
	case bool:
		return value.(bool)
	case string:
		if value == "true" {
			return true
		}
		if value == "false" {
			return false
		}
		return defaultValue
	default:
		return defaultValue
	}
}
