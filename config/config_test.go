package config

import (
	"errors"
	"testing"
	"time"

	"github.com/fern4lvarez/piladb/config/vars"
	"github.com/fern4lvarez/piladb/pkg/uuid"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	if config.Values == nil {
		t.Fatal(errors.New("Config is nil"))
	}

	inputOutput := []struct {
		input, output interface{}
	}{
		{config.Values.Name, CONFIG},
		{config.Values.ID, uuid.New(CONFIG)},
		{len(config.Values.Stacks), 0},
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
		{vars.MaxStackSize, -1},
	}
	for _, io := range inputOutput {
		if value := config.Get(io.input); value != io.output {
			t.Errorf("Values is %v, expected %v", value, io.output)
		}
	}
}

func TestConfigGet(t *testing.T) {
	config := NewConfig()

	stackID := config.Values.CreateStack("foo", time.Now())
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
