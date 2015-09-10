package pila

import "testing"

func TestNewStack(t *testing.T) {
	stack := NewStack("test-stack")

	if stack == nil {
		t.Fatal("stack is nil")
	}

	if stack.ID.String() != "2f44edeaa249ba81db20e9ddf000ba65" {
		t.Errorf("stack.ID is %s, expected %s", stack.ID.String(), "2f44edeaa249ba81db20e9ddf000ba65")
	}
	if stack.Name != "test-stack" {
		t.Errorf("stack.Name is %s, expected %s", stack.Name, "test-stack")
	}
	if stack.Database != nil {
		t.Error("stack.Database is not nil")
	}
	if stack.Base == nil {
		t.Fatalf("stack.Base is nil")
	}
	if stack.Base.Size() != 0 {
		t.Fatalf("stack.Base.Size() is %d, expected %d", stack.Base.Size(), 0)
	}

}
