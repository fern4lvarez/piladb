package config

import (
	"testing"

	"github.com/fern4lvarez/piladb/config/vars"
)

func TestMaxStackSize(t *testing.T) {
	c := NewConfig()

	inputOutput := []struct {
		input  interface{}
		output int
	}{
		{8, 8},
		{-1, -1},
		{23.7, 23},
		{"3", 3},
		{"foo", -1},
		{[]byte("foo"), -1},
	}

	for _, io := range inputOutput {
		c.Set(vars.MaxStackSize, io.input)

		if s := c.MaxStackSize(); s != io.output {
			t.Errorf("MaxStackSize is %d, expected %d", s, io.output)
		}
	}
}
