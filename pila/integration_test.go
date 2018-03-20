package pila

import (
	"reflect"
	"testing"
	"time"
)

func TestIntegrationBasic(t *testing.T) {
	pila := NewPila()

	db := NewDatabase("db")
	pila.AddDatabase(db)

	stack1 := NewStack("stack1", time.Now())
	stack2 := NewStack("stack2", time.Now())
	db.AddStack(stack1)
	db.AddStack(stack2)

	elements := []interface{}{"foo", "bar", "baz"}

	for _, element := range elements {
		stack1.Push(element)
	}

	if stack1.Size() != 3 {
		t.Errorf("stack1.Size is %d, expected %d", stack1.Size(), 3)
	}
	if stack2.Size() != 0 {
		t.Errorf("stack2.Size is %d, expected %d", stack2.Size(), 0)
	}

	if stack1.Peek() != interface{}("baz") {
		t.Errorf("stack1.Peek is %v, expected %v", stack1.Peek(), interface{}("baz"))
	}
	if stack2.Peek() != nil {
		t.Errorf("stack2.Peek is %v, expected %v", stack2.Peek(), nil)
	}

	for {
		element, _ := stack1.Pop()

		if element != nil {
			stack2.Push(element)
		} else {
			break
		}
	}

	if stack1.Size() != 0 {
		t.Errorf("stack1.Size is %d, expected %d", stack1.Size(), 0)
	}
	if stack2.Size() != 3 {
		t.Errorf("stack2.Size is %d, expected %d", stack2.Size(), 3)
	}

	if stack1.Peek() != nil {
		t.Errorf("stack1.Peek is %v, expected %v", stack1.Peek(), nil)
	}
	if stack2.Peek() != interface{}("foo") {
		t.Errorf("stack2.Peek is %v, expected %v", stack2.Peek(), interface{}("foo"))
	}

	if ok := db.RemoveStack(stack1.ID); !ok {
		t.Errorf("database %s failed on removing stack %s", db.Name, stack1.Name)
	}
	if _, ok := db.Stacks[stack1.ID]; ok {
		t.Errorf("stack1 %s was found in database %s", stack1.Name, db.Name)
	}

	stack2Copy, ok := db.Stacks[stack2.ID]
	if !ok {
		t.Errorf("stack2 %v was not found in database %s", stack2.Name, db.Name)
	}

	if !reflect.DeepEqual(&stack2Copy.Database, &db) {
		t.Errorf("stack2Copy.Database is %v, expected %v", &stack2Copy.Database, &db)
	}
	if !reflect.DeepEqual(&stack2Copy.Database.Pila, &db.Pila) {
		t.Errorf("stack2Copy.Database.Pila is %v, expected %v", &stack2Copy.Database.Pila, &db.Pila)
	}
}
