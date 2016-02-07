package pila

import (
	"fmt"
	"testing"
)

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
	if stack.base == nil {
		t.Fatalf("stack.base is nil")
	}
	if stack.base.Size() != 0 {
		t.Fatalf("stack.base.Size() is %d, expected %d", stack.base.Size(), 0)
	}
}

func TestStackPush(t *testing.T) {
	stack := NewStack("test-stack")
	stack.Push(1)

	if stack.Size() != 1 {
		t.Errorf("stack.Size() is %d, expected %d", stack.base.Size(), 0)
	}

	stack.Push(2)
	stack.Push(struct{ id string }{id: "test"})

	if stack.Size() != 3 {
		t.Errorf("stack.base.Size() is %d, expected %d", stack.base.Size(), 3)
	}
}

func TestStackPop(t *testing.T) {
	stack := NewStack("test-stack")
	stack.Push("test")
	stack.Push(8)

	element, ok := stack.Pop()
	if !ok {
		t.Errorf("stack.Pop() not ok")
	}
	if element != 8 {
		t.Errorf("element is %v, expected %v", element, 8)
	}
	if stack.Size() != 1 {
		t.Errorf("stack.Size() is %d, expected %d", stack.Size(), 1)
	}
}

func TestStackPop_False(t *testing.T) {
	stack := NewStack("test-stack")
	_, ok := stack.Pop()
	if ok {
		t.Error("stack.Pop() is ok")
	}
}

func TestStackSize(t *testing.T) {
	stack := NewStack("test-stack")
	if stack.Size() != 0 {
		t.Errorf("stack.Size() is %d, expected %d", stack.Size(), 0)
	}

	for i := 0; i < 15; i++ {
		stack.Push(i)
	}
	if stack.Size() != 15 {
		t.Errorf("stack.Size() is %d, expected %d", stack.Size(), 15)
	}

	for i := 0; i < 5; i++ {
		stack.Pop()
	}
	if stack.Size() != 10 {
		t.Errorf("stack.Size() is %d, expected %d", stack.Size(), 10)
	}
}

func TestStackPeek(t *testing.T) {
	stack := NewStack("test-stack")
	stack.Push("test")
	stack.Push(8)

	element := stack.Peek()
	if element != 8 {
		t.Errorf("element is %v, expected %v", element, 8)
	}
}

func TestStackStatus(t *testing.T) {
	stack := NewStack("test-stack")
	stack.Push("test")
	stack.Push(8)
	stack.Push(5.87)
	stack.Push([]byte("test"))

	expectedStatus := fmt.Sprintf(`{"id":"2f44edeaa249ba81db20e9ddf000ba65","name":"test-stack","peek":"dGVzdA==","size":4}`)
	if status, err := stack.Status(); err != nil {
		t.Fatal(err)
	} else if string(status) != expectedStatus {
		t.Errorf("status is %s, expected %s", string(status), expectedStatus)
	}
}

func TestStackStatus_Empty(t *testing.T) {
	stack := NewStack("test-stack")

	expectedStatus := fmt.Sprintf(`{"id":"2f44edeaa249ba81db20e9ddf000ba65","name":"test-stack","peek":null,"size":0}`)
	if status, err := stack.Status(); err != nil {
		t.Fatal(err)
	} else if string(status) != expectedStatus {
		t.Errorf("status is %s, expected %s", string(status), expectedStatus)
	}
}

func TestStackStatus_Error(t *testing.T) {
	// From https://golang.org/src/encoding/json/encode.go?s=5438:5481#L125
	// Channel, complex, and function values cannot be encoded in JSON.
	// Attempting to encode such a value causes Marshal to return
	// an UnsupportedTypeError.

	ch := make(chan int)

	stack := NewStack("test-stack")
	stack.Push(ch)

	if _, err := stack.Status(); err == nil {
		t.Error("err is nil, expected UnsupportedTypeError")
	}
}
