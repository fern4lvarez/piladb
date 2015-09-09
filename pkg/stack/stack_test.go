package stack

import "testing"

func TestNewStack(t *testing.T) {
	stack := NewStack()
	if stack.head != nil {
		t.Error("stack.head is not nil")
	}
	if stack.size != 0 {
		t.Errorf("stack.size is %v, expected 0", stack.size)
	}
}

func TestStackPush(t *testing.T) {
	stack := NewStack()
	stack.Push(8)

	if stack.head == nil {
		t.Error("stack.head is nil")
	} else if stack.head.data != 8 {
		t.Errorf("stack.head data is %v, expected 8", stack.head.data)
	}
	if stack.size != 1 {
		t.Errorf("stack.size is %v, expected %v", stack.size, 1)
	}
}

func TestStackPop(t *testing.T) {
	stack := NewStack()
	stack.Push("test")
	stack.Push(8)

	element, ok := stack.Pop()
	if !ok {
		t.Errorf("stack.Pop() not ok")
	}
	if element != 8 {
		t.Errorf("element is %v, expected %v", element, 8)
	}
	if stack.head == nil {
		t.Error("stack.head is nil")
	} else if stack.head.data != "test" {
		t.Errorf("stack.head data is %v, expected %v", stack.head.data, "test")
	}
	if stack.size != 1 {
		t.Errorf("stack.size is %v, expected %v", stack.size, 1)
	}
}

func TestStackPop_False(t *testing.T) {
	stack := NewStack()
	_, ok := stack.Pop()
	if ok {
		t.Error("stack.Pop() is ok")
	}
}

func TestStackSize(t *testing.T) {
	stack := NewStack()
	if stack.Size() != 0 {
		t.Errorf("stack.Size() is %v, expected %v", stack.Size(), 0)
	}

	for i := 0; i < 15; i++ {
		stack.Push(i)
	}
	if stack.Size() != 15 {
		t.Errorf("stack.Size() is %v, expected %v", stack.Size(), 15)
	}

	for i := 0; i < 5; i++ {
		stack.Pop()
	}
	if stack.Size() != 10 {
		t.Errorf("stack.Size() is %v, expected %v", stack.Size(), 10)
	}
}

func TestStackPeek(t *testing.T) {
	stack := NewStack()
	if stack.Peek() != nil {
		t.Error("stack.Peek() is not nil")
	}

	stack.Push(9.35)
	if stack.Peek() != 9.35 {
		t.Errorf("stack.Peek() is %v, expected %v", stack.Peek(), 9.35)
	}

	stack.Push(true)
	if stack.Peek() == false {
		t.Errorf("stack.Peek() is %v, expected %v", stack.Peek(), true)
	}

	stack.Push("one")
	stack.Push("two")
	stack.Push("three")
	stack.Pop()
	if stack.Peek() != "two" {
		t.Errorf("stack.Peek() is %v, expected %v", stack.Peek(), "two")
	}

}
