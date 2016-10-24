package config

import (
	"errors"
	"testing"

	"github.com/fern4lvarez/piladb/config/vars"
	"github.com/fern4lvarez/piladb/pkg/uuid"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	if config.Values == nil {
		t.Fatal(errors.New("Config is nil"))
	}

	inputOutput := []struct {
		input  interface{}
		output interface{}
	}{
		{input: config.Values.Name, output: CONFIG},
		{input: config.Values.ID, output: uuid.New(CONFIG)},
		{input: len(config.Values.Stacks), output: 0},
	}

	for _, io := range inputOutput {
		if io.input != io.output {
			t.Errorf("got %v, expected %v", io.input, io.output)
		}
	}
}

func TestConfigDefault(t *testing.T) {
	config := NewConfig().Default()
	inputOutput := []struct {
		input  string
		output interface{}
	}{
		{input: vars.MaxSizeOfStack, output: 5},
	}
	for _, io := range inputOutput {
		if value := config.Get(io.input); value != io.output {
			t.Errorf("Values is %v, expected %v", value, io.output)
		}
	}
}

func TestConfigGet(t *testing.T) {
	config := NewConfig()

	stackID := config.Values.CreateStack("foo")
	s, _ := config.Values.Stacks[stackID]
	s.Push("bar")
	expectedValue := s.Peek()

	if value := config.Get("foo"); value != expectedValue {
		t.Errorf("Values is %s, expected %s", value, expectedValue)
	}

	if value := config.Get("no-exist"); value != nil {
		t.Errorf("Values is %s, expected nil", value)
	}
}

func TestConfigSet(t *testing.T) {
	config := NewConfig()
	expectedValues := []string{"bar", "baz", "bam"}

	for _, expectedValue := range expectedValues {
		config.Set("foo", expectedValue)
		s, _ := config.Values.Stacks[uuid.New(CONFIG+"foo")]
		if value := s.Peek(); value != expectedValue {
			t.Errorf("Values is %s, expected %s", value, expectedValue)
		}
	}
}
