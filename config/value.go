package config

import (
	"strconv"

	"github.com/fern4lvarez/piladb/config/vars"
)

// MaxStackSize returns the value of MAX_STACK_SIZE.
// Type: int, Default: -1
func (c *Config) MaxStackSize() int {
	maxSize := c.Get(vars.MaxStackSize)
	return intValue(maxSize, -1)
}

// ReadTimeout returns the value of READ_TIMEOUT.
// Type: int, Default: 30
func (c *Config) ReadTimeout() int {
	readTimeout := c.Get(vars.ReadTimeout)
	return intValue(readTimeout, 30)
}

// WriteTimeout returns the value of WRITE_TIMEOUT.
// Type: int, Default: 45
func (c *Config) WriteTimeout() int {
	writeTimeout := c.Get(vars.WriteTimeout)
	return intValue(writeTimeout, 45)
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
