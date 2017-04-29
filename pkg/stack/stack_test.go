package stack

import "testing"

func TestNewStack(t *testing.T) {
	stack := NewStack()
	if stack.head != nil {
		t.Error("stack.head is not nil")
	}
	if stack.tail != nil {
		t.Error("stack.tail is not nil")
	}
	if stack.size != 0 {
		t.Errorf("stack.size is %v, expected 0", stack.size)
	}
}

func TestStackPush(t *testing.T) {
	stack := NewStack()
	stack.Push(8)

	if stack.head == nil {
		t.Fatal("stack.head is nil")
	}
	if stack.head.data != 8 {
		t.Errorf("stack.head data is %v, expected 8", stack.head.data)
	}
	if stack.head.down != nil {
		t.Errorf("stack.head.down is %v, expected nil", stack.head.down)
	}
	if stack.tail == nil {
		t.Fatal("stack.tail is nil")
	}
	if stack.tail.data != 8 {
		t.Errorf("stack.tail data is %v, expected 8", stack.tail.data)
	}
	if stack.tail.down != nil {
		t.Errorf("stack.tail.down is %v, expected nil", stack.tail.down)
	}
	if stack.size != 1 {
		t.Errorf("stack.size is %v, expected %v", stack.size, 1)
	}
}

func TestStackPush_TwoElements(t *testing.T) {
	stack := NewStack()
	stack.Push(8)
	expectedNext := stack.head
	stack.Push("test")

	if stack.head == nil {
		t.Fatal("stack.head is nil")
	}
	if stack.head.data != "test" {
		t.Errorf("stack.head data is %v, expected test", stack.head.data)
	}
	if stack.head.down != expectedNext {
		t.Errorf("stack.head.down is %v, expected %v", stack.head.down, expectedNext)
	}
	if stack.head.up != nil {
		t.Errorf("stack.head.up is %v, expected nil", stack.head.up)
	}
	if stack.tail == nil {
		t.Fatal("stack.tail is nil")
	}
	if stack.tail.data != 8 {
		t.Errorf("stack.tail data is %v, expected 8", stack.tail.data)
	}
	if stack.tail.down != nil {
		t.Errorf("stack.tail.down is %v, expected nil", stack.tail.down)
	}
	if stack.tail.up != stack.head {
		t.Errorf("stack.head.up is %v, expected %v", stack.head.up, stack.head)
	}
	if stack.size != 2 {
		t.Errorf("stack.size is %v, expected %v", stack.size, 2)
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
		t.Fatal("stack.head is nil")
	}
	if stack.head.data != "test" {
		t.Errorf("stack.head data is %v, expected %v", stack.head.data, "test")
	}
	if stack.head.up != nil {
		t.Errorf("stack.head.up data is %v, expected %v", stack.head.up, nil)
	}
	if stack.tail == nil {
		t.Fatal("stack.tail is nil")
	}
	if stack.tail.data != "test" {
		t.Errorf("stack.tail.data is %v, expected %v", stack.tail.data, "test")
	}
	if stack.tail.down != nil {
		t.Errorf("stack.tail.data is %v, expected nil", stack.tail.data)
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

func TestStackBase(t *testing.T) {
	stack := NewStack()
	stack.Base("test")

	if stack.tail == nil {
		t.Fatal("stack.tail is nil")
	}
	if stack.tail.data != "test" {
		t.Errorf("stack.tail data is %v, expected %v", stack.tail.data, "data")
	}
	if stack.tail.down != nil {
		t.Errorf("stack.tail.down is %v, expected nil", stack.tail.down)
	}
	if stack.tail.up != nil {
		t.Errorf("stack.tail.up is %v, expected nil", stack.tail)
	}
	if stack.head == nil {
		t.Fatal("stack.head is nil")
	}
	if stack.head.data != "test" {
		t.Errorf("stack.head.data is %v, expected %v", stack.head.data, "test")
	}
	if stack.head.down != nil {
		t.Errorf("stack.head.down is %v, expected nil", stack.head.down)
	}
	if stack.size != 1 {
		t.Errorf("stack.size is %v, expected %v", stack.size, 1)
	}
}

func TestStackBase_TwoElements(t *testing.T) {
	stack := NewStack()
	stack.Base("test")
	stack.Base(8)

	if stack.tail == nil {
		t.Fatal("stack.tail is nil")
	}
	if stack.tail.data != 8 {
		t.Errorf("stack.tail data is %v, expected %v", stack.tail.data, 8)
	}
	if stack.tail.down != nil {
		t.Errorf("stack.tail.down is %v, expected nil", stack.tail.down)
	}
	if stack.tail.up.data != "test" {
		t.Errorf("stack.tail.up is %v, expected %v", stack.tail.up.data, "test")
	}
	if stack.head == nil {
		t.Fatal("stack.head is nil")
	}
	if stack.head.data != "test" {
		t.Errorf("stack.head.data is %v, expected %v", stack.head.data, "test")
	}
	if stack.head.down.data != 8 {
		t.Errorf("stack.head.down is %v, expected %v", stack.head.down.data, 8)
	}
	if stack.size != 2 {
		t.Errorf("stack.size is %v, expected %v", stack.size, 2)
	}
}

func TestStackBase_MoreElements(t *testing.T) {
	stack := NewStack()
	stack.Base("test")
	stack.Base(8)
	stack.Base(true)

	if stack.tail == nil {
		t.Fatal("stack.tail is nil")
	}
	if stack.tail.data != true {
		t.Errorf("stack.tail data is %v, expected %v", stack.tail.data, true)
	}
	if stack.tail.down != nil {
		t.Errorf("stack.tail.down is %v, expected nil", stack.tail.down)
	}
	if stack.tail.up.data != 8 {
		t.Errorf("stack.tail.up is %v, expected %v", stack.tail.up.data, 8)
	}
	if stack.head == nil {
		t.Fatal("stack.head is nil")
	}
	if stack.head.data != "test" {
		t.Errorf("stack.head.data is %v, expected %v", stack.head.data, "test")
	}
	if stack.head.down.data != 8 {
		t.Errorf("stack.head.down is %v, expected %v", stack.head.down.data, 8)
	}
	if stack.size != 3 {
		t.Errorf("stack.size is %v, expected %v", stack.size, 3)
	}
}

func TestStackBase_NoEmpty(t *testing.T) {
	stack := NewStack()
	stack.Push("test")
	stack.Base(8)

	if stack.tail == nil {
		t.Fatal("stack.tail is nil")
	}
	if stack.tail.data != 8 {
		t.Errorf("stack.tail data is %v, expected %v", stack.tail.data, 8)
	}
	if stack.tail.down != nil {
		t.Errorf("stack.tail.down is %v, expected nil", stack.tail.down)
	}
	if stack.tail.up.data != "test" {
		t.Errorf("stack.tail.up is %v, expected %v", stack.tail.up.data, "test")
	}
	if stack.head == nil {
		t.Fatal("stack.head is nil")
	}
	if stack.head.data != "test" {
		t.Errorf("stack.head.data is %v, expected %v", stack.head.data, "test")
	}
	if stack.head.down.data != 8 {
		t.Errorf("stack.head.down is %v, expected %v", stack.head.down.data, 8)
	}
	if stack.size != 2 {
		t.Errorf("stack.size is %v, expected %v", stack.size, 2)
	}
}

func TestStackBase_NoEmptyMore(t *testing.T) {
	stack := NewStack()
	stack.Push("test")
	stack.Push(3.14)
	stack.Push(false)
	stack.Base("foo")
	stack.Base(8)

	if stack.tail == nil {
		t.Fatal("stack.tail is nil")
	}
	if stack.tail.data != 8 {
		t.Errorf("stack.tail data is %v, expected %v", stack.tail.data, 8)
	}
	if stack.tail.down != nil {
		t.Errorf("stack.tail.down is %v, expected nil", stack.tail.down)
	}
	if stack.tail.up.data != "foo" {
		t.Errorf("stack.tail.up is %v, expected %v", stack.tail.up.data, "foo")
	}
	if stack.head == nil {
		t.Fatal("stack.head is nil")
	}
	if stack.head.data != false {
		t.Errorf("stack.head.data is %v, expected %v", stack.head.data, false)
	}
	if stack.head.down.data != 3.14 {
		t.Errorf("stack.head.down is %v, expected %v", stack.head.down.data, 3.14)
	}
	if stack.size != 5 {
		t.Errorf("stack.size is %v, expected %v", stack.size, 5)
	}
}

func TestStackSweep(t *testing.T) {
	stack := NewStack()
	stack.Push("test")
	stack.Push(8)

	element, ok := stack.Sweep()
	if !ok {
		t.Errorf("stack.Sweep() not ok")
	}
	if element != "test" {
		t.Errorf("element is %v, expected %v", element, "test")
	}
	if stack.tail == nil {
		t.Fatal("stack.tail is nil")
	}
	if stack.tail.data != 8 {
		t.Errorf("stack.tail data is %v, expected %v", stack.tail.data, 8)
	}
	if stack.tail.down != nil {
		t.Errorf("stack.tail.down is %v, expected nil", stack.tail.down)
	}
	if stack.tail.up != nil {
		t.Errorf("stack.tail.up is %v, expected nil", stack.tail)
	}
	if stack.head == nil {
		t.Fatal("stack.head is nil")
	}
	if stack.head.data != 8 {
		t.Errorf("stack.head.data is %v, expected %v", stack.head.data, 8)
	}
	if stack.head.down != nil {
		t.Errorf("stack.head.down is %v, expected nil", stack.head.down)
	}
	if stack.size != 1 {
		t.Errorf("stack.size is %v, expected %v", stack.size, 1)
	}
}

func TestStackSweep_OneElement(t *testing.T) {
	stack := NewStack()
	stack.Push("test")

	element, ok := stack.Sweep()
	if !ok {
		t.Errorf("stack.Sweep() not ok")
	}
	if element != "test" {
		t.Errorf("element is %v, expected %v", element, "test")
	}
	if stack.tail != nil {
		t.Errorf("stack.tail is %v, expected nil", stack.tail)
	}
	if stack.head != nil {
		t.Errorf("stack.tail is %v, expected nil", stack.head)
	}
	if stack.size != 0 {
		t.Errorf("stack.size is %v, expected %v", stack.size, 0)
	}
}

func TestStackSweep_More(t *testing.T) {
	stack := NewStack()
	stack.Push("test")
	stack.Push(8)
	stack.Push(10)

	element, ok := stack.Sweep()
	if !ok {
		t.Errorf("stack.Sweep() not ok")
	}
	if element != "test" {
		t.Errorf("element is %v, expected %v", element, "test")
	}
	if stack.tail == nil {
		t.Fatal("stack.tail is nil")
	}
	if stack.tail.data != 8 {
		t.Errorf("stack.tail data is %v, expected %v", stack.tail.data, 8)
	}
	if stack.head == nil {
		t.Fatal("stack.head is nil")
	}
	if stack.head.data != 10 {
		t.Errorf("stack.head.data is %v, expected %v", stack.head.data, 10)
	}
	if stack.head.down != stack.tail {
		t.Errorf("stack.head.down is %v, expected %v", stack.head.down, stack.tail)
	}
	if stack.size != 2 {
		t.Errorf("stack.size is %v, expected %v", stack.size, 2)
	}

	element, ok = stack.Sweep()
	if !ok {
		t.Errorf("stack.Sweep() not ok")
	}
	if element != 8 {
		t.Errorf("element is %v, expected %v", element, 8)
	}
	if stack.tail == nil {
		t.Fatal("stack.tail is nil")
	}
	if stack.tail.data != 10 {
		t.Errorf("stack.tail data is %v, expected %v", stack.tail.data, 8)
	}
	if stack.head == nil {
		t.Fatal("stack.head is nil")
	}
	if stack.head.data != 10 {
		t.Errorf("stack.head.data is %v, expected %v", stack.head.data, 10)
	}
	if stack.size != 1 {
		t.Errorf("stack.size is %v, expected %v", stack.size, 1)
	}
}

func TestStackSweep_AndPush(t *testing.T) {
	stack := NewStack()
	stack.Push("test")
	stack.Push(8)
	stack.Push(10)

	_, _ = stack.Sweep()
	stack.Push("foo")

	_, _ = stack.Sweep()
	stack.Push(0)

	element, ok := stack.Sweep()
	stack.Push(true)
	if !ok {
		t.Errorf("stack.Sweep() not ok")
	}
	if element != 10 {
		t.Errorf("element is %v, expected %v", element, 10)
	}
	if stack.tail == nil {
		t.Fatal("stack.tail is nil")
	}
	if stack.tail.data != "foo" {
		t.Errorf("stack.tail data is %v, expected %v", stack.tail.data, "foo")
	}
	if stack.tail.down != nil {
		t.Errorf("stack.tail.down is %v, expected nil", stack.tail.down)
	}
	if stack.head == nil {
		t.Fatal("stack.head is nil")
	}
	if stack.head.data != true {
		t.Errorf("stack.head.data is %v, expected %v", stack.head.data, true)
	}
	if stack.size != 3 {
		t.Errorf("stack.size is %v, expected %v", stack.size, 3)
	}
}

func TestStackSweep_False(t *testing.T) {
	stack := NewStack()
	_, ok := stack.Sweep()
	if ok {
		t.Error("stack.Sweep() is ok")
	}
}

func TestStackSweepPush(t *testing.T) {
	stack := NewStack()
	stack.Push("test")
	stack.Push(8)

	element, ok := stack.SweepPush("foo")
	if !ok {
		t.Errorf("stack.SweepPush(foo) not ok")
	}
	if element != "test" {
		t.Errorf("element is %v, expected %v", element, "test")
	}
	if stack.tail == nil {
		t.Fatal("stack.tail is nil")
	}
	if stack.tail.data != 8 {
		t.Errorf("stack.tail data is %v, expected %v", stack.tail.data, 8)
	}
	if stack.tail.down != nil {
		t.Errorf("stack.tail.down is %v, expected nil", stack.tail.down)
	}
	if stack.tail.up.data != "foo" {
		t.Errorf("stack.tail.up.data is %v, expected %v", stack.tail.up.data, "foo")
	}
	if stack.head == nil {
		t.Fatal("stack.head is nil")
	}
	if stack.head.data != "foo" {
		t.Errorf("stack.head.data is %v, expected %v", stack.head.data, "foo")
	}
	if stack.head.down.data != 8 {
		t.Errorf("stack.head.down.data is %v, expected %v", stack.head.down.data, 8)
	}
	if stack.size != 2 {
		t.Errorf("stack.size is %v, expected %v", stack.size, 2)
	}
}

func TestStackSweepPush_More(t *testing.T) {
	stack := NewStack()
	stack.Push(8)
	stack.Push("{'a':'b'}")
	stack.Push(23.34)
	element, ok := stack.SweepPush("foo")
	if !ok {
		t.Errorf("stack.SweepPush(foo) not ok")
	}
	if element != 8 {
		t.Errorf("element is %v, expected %v", element, 8)
	}

	if stack.head == nil {
		t.Fatal("stack.head is nil")
	}
	if stack.head.data != "foo" {
		t.Errorf("stack.head.data is %v, expected %v", stack.head.data, "foo")
	}
	if stack.head.up != nil {
		t.Errorf("stack.head.up is %v, expected nil", stack.head.up)
	}
	if stack.head.down.data != 23.34 {
		t.Errorf("stack.head.down.data is %v, expected %v", stack.head.down, 23.34)
	}
	if stack.tail == nil {
		t.Fatal("stack.tail is nil")
	}
	if stack.tail.data != "{'a':'b'}" {
		t.Errorf("stack.tail data is %v, expected %v", stack.tail.data, "{'a':'b'}")
	}
	if stack.tail.down != nil {
		t.Errorf("stack.tail.down is %v, expected nil", stack.tail.down)
	}
	if stack.size != 3 {
		t.Errorf("stack.size is %v, expected %v", stack.size, 3)
	}
}

func TestStackSweepPush_OneElement(t *testing.T) {
	stack := NewStack()
	stack.Push(8)
	element, ok := stack.SweepPush("foo")
	if !ok {
		t.Errorf("stack.SweepPush(foo) not ok")
	}
	if element != 8 {
		t.Errorf("element is %v, expected %v", element, 8)
	}

	if stack.head == nil {
		t.Fatal("stack.head is nil")
	}
	if stack.head.data != "foo" {
		t.Errorf("stack.head.data is %v, expected %v", stack.head.data, "foo")
	}
	if stack.head.down != nil {
		t.Errorf("stack.head.down is %v, expected nil", stack.head.down)
	}
	if stack.tail == nil {
		t.Fatal("stack.tail is nil")
	}
	if stack.tail.data != "foo" {
		t.Errorf("stack.tail data is %v, expected %v", stack.tail.data, "foo")
	}
	if stack.tail.down != nil {
		t.Errorf("stack.tail.down is %v, expected nil", stack.tail.down)
	}
	if stack.size != 1 {
		t.Errorf("stack.size is %v, expected %v", stack.size, 1)
	}
}

func TestStackSweepPush_False(t *testing.T) {
	stack := NewStack()
	_, ok := stack.SweepPush(8)
	if ok {
		t.Error("stack.SweepPush(8) is ok")
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
	if !stack.Peek().(bool) {
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

func TestStackFlush(t *testing.T) {
	stack := NewStack()

	stack.Push("one")
	stack.Push("two")
	stack.Push("three")

	stack.Flush()

	if stack.Peek() != nil {
		t.Error("stack.Peek() is not nil")
	}
	if stack.Size() != 0 {
		t.Errorf("stack is not empty")
	}
}

func TestStackRace(t *testing.T) {
	stack := NewStack()
	go func() { stack.Push(1) }()
	go func() { stack.Pop() }()
	go func() { stack.Size() }()
	go func() { stack.Peek() }()
	go func() { stack.Flush() }()
}
