// Package config implements the configuration management of piladb.
package config

import (
	"time"

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
		sID := c.Values.CreateStack(key, time.Now().UTC())
		s, _ = c.Values.Stacks[sID]
	}

	s.Push(value)
}
