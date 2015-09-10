package pila

import "testing"

func TestNewDatabase(t *testing.T) {
	db := NewDatabase("test-1")

	if db.ID.String() != "2b87e5d8b7d3d853514c8d0801fbcf46" {
		t.Errorf("db.ID is %v, expected %v", db.ID, "2b87e5d8b7d3d853514c8d0801fbcf46")
	}
	if db.Name != "test-1" {
		t.Errorf("db.Name is %v, expected %v", db.Name, "test-1")
	}
	if db.Pila != nil {
		t.Error("db.Pila is not nil")
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
