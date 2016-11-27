package config

import (
	"testing"
	"time"

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

func TestReadTimeout(t *testing.T) {
	c := NewConfig()

	inputOutput := []struct {
		input  interface{}
		output time.Duration
	}{
		{8, 8},
		{-1, 30},
		{23.7, 23},
		{"3", 3},
		{"foo", 30},
		{[]byte("foo"), 30},
	}

	for _, io := range inputOutput {
		c.Set(vars.ReadTimeout, io.input)

		if s := c.ReadTimeout(); s != io.output {
			t.Errorf("ReadTimeout is %d, expected %d", s, io.output)
		}
	}
}

func TestWriteTimeout(t *testing.T) {
	c := NewConfig()

	inputOutput := []struct {
		input  interface{}
		output time.Duration
	}{
		{8, 8},
		{-1, 45},
		{23.7, 23},
		{"3", 3},
		{"foo", 45},
		{[]byte("foo"), 45},
	}

	for _, io := range inputOutput {
		c.Set(vars.WriteTimeout, io.input)

		if s := c.WriteTimeout(); s != io.output {
			t.Errorf("WriteTimeout is %d, expected %d", s, io.output)
		}
	}
}

func TestPort(t *testing.T) {
	c := NewConfig()

	inputOutput := []struct {
		input  interface{}
		output int
	}{
		{8090, 8090},
		{-1, 1205},
		{23676.7, 23676},
		{"3756", 3756},
		{"foo", 1205},
		{[]byte("foo"), 1205},
		{6736373635, 1205},
	}

	for _, io := range inputOutput {
		c.Set(vars.Port, io.input)

		if s := c.Port(); s != io.output {
			t.Errorf("Port is %d, expected %d", s, io.output)
		}
	}
}
