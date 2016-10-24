package config

import (
	"github.com/fern4lvarez/piladb/config/vars"
	"github.com/fern4lvarez/piladb/pila"
	"github.com/fern4lvarez/piladb/pkg/uuid"
)

// CONFIG represents the name of the database that will hold
// all the config values.
const CONFIG = "_config"

// Config represents a Database containing all
// configuration values that will be
// updated and consumed by piladb.
type Config struct {
	Values *pila.Database
}

// NewConfig creates a new Config with empty values.
func NewConfig() *Config {
	return &Config{Values: pila.NewDatabase(CONFIG)}
}

// Default sets the default values to the Config.
func (c *Config) Default() *Config {
	// Infinite size
	c.Set(vars.MaxSizeOfStack, 5)

	return c
}

// Get gets a config value from a key.
func (c *Config) Get(key string) interface{} {
	s, ok := c.Values.Stacks[uuid.New(CONFIG+key)]
	if !ok {
		return nil
	}
	return s.Peek()
}

// Set sets a config value having a key and the value.
func (c *Config) Set(key string, value interface{}) {
	s, ok := c.Values.Stacks[uuid.New(CONFIG+key)]
	if !ok {
		sID := c.Values.CreateStack(key)
		s, _ = c.Values.Stacks[sID]
	}

	s.Push(value)
}
