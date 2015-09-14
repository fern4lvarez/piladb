package pila

import (
	"reflect"
	"testing"
)

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

func TestDatabaseCreateStack(t *testing.T) {
	db := NewDatabase("test-db")
	id := db.CreateStack("test-stack")

	if id == nil {
		t.Fatal("stack ID is nil")
	}

	stack, ok := db.Stacks[id]
	if !ok {
		t.Fatal("stack not found in database")
	}
	if stack.ID != id {
		t.Errorf("stack ID is %v, expected %v", stack.ID, id)
	}
	if stack.Name != "test-stack" {
		t.Errorf("stack Name is %s, expected %s", stack.Name, "test-stack")
	}

	if !reflect.DeepEqual(stack.Database, db) {
		t.Errorf("stack Database is %v, expected %v", stack.Database, db)
	}
}

func TestDatabaseAddStack(t *testing.T) {
	db := NewDatabase("test-db")
	stack := NewStack("test-stack")

	err := db.AddStack(stack)
	if err != nil {
		t.Fatal("err is not nil")
	}

	stack2, ok := db.Stacks[stack.ID]
	if !ok {
		t.Error("Stack not found in Database")
	}
	if !reflect.DeepEqual(stack2, stack) {
		t.Errorf("Stack is %v, expected %v", stack2, stack)
	}
	if !reflect.DeepEqual(stack.Database, db) {
		t.Errorf("Stack.Database is %v, expected %v", stack.Database, db)
	}
}

func TestDatabaseAddStack_ErrorStackAlreadyAdded(t *testing.T) {
	db := NewDatabase("test-db")
	db2 := NewDatabase("test-db-2")
	stack := NewStack("test-stack")

	err := db.AddStack(stack)
	if err != nil {
		t.Fatal("err is not nil")
	}

	err = db2.AddStack(stack)
	if err == nil {
		t.Fatal("err is nil")
	}
}

func TestDatabaseAddStack_ErrorDBAlreadyContainsStack(t *testing.T) {
	db := NewDatabase("test-db")
	stack := NewStack("test-stack")
	stack2 := NewStack("test-stack")

	err := db.AddStack(stack)
	if err != nil {
		t.Fatal("err is not nil")
	}

	err = db.AddStack(stack2)
	if err == nil {
		t.Fatal("err is nil")
	}
}
