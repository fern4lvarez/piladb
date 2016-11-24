package config

import (
	"strconv"

	"github.com/fern4lvarez/piladb/config/vars"
)

// MaxStackSize returns the value of MAX_STACK_SIZE.
// Type: int, Default: -1
func (c *Config) MaxStackSize() int {
	maxSize := c.Get(vars.MaxStackSize)
	switch maxSize.(type) {
	case int:
		return maxSize.(int)
	case float64:
		return int(maxSize.(float64))
	case string:
		i, err := strconv.Atoi(maxSize.(string))
		if err != nil {
			return -1
		}
		return i
	default:
		return -1
	}
}
