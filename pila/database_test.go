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

func TestDatabaseRemoveStack(t *testing.T) {
	db := NewDatabase("test-db")
	stack := NewStack("test-stack")

	err := db.AddStack(stack)
	if err != nil {
		t.Fatal("err is not nil")
	}

	ok := db.RemoveStack(stack.ID)
	if !ok {
		t.Errorf("stack %s was not removed from database %s", stack.Name, db.Name)
	}

	_, ok = db.Stacks[stack.ID]
	if ok {
		t.Errorf("stack %s was found in database %s", stack.Name, db.Name)
	}

	if stack.Database != nil {
		t.Errorf("stack %s still associated to database %s", stack.Name, stack.Database.Name)
	}

	if stack.base != nil {
		t.Errorf("stack %s still points to a base stack", stack.Name)
	}

}

func TestDatabaseRemoveStack_False(t *testing.T) {
	db := NewDatabase("test-db")
	stack := NewStack("test-stack")

	ok := db.RemoveStack(stack.ID)
	if ok {
		t.Errorf("stack %v was removed from database %v", stack.Name, db.Name)
	}
}

func TestDatabaseStatus(t *testing.T) {
	db := NewDatabase("db")
	s0ID := db.CreateStack("s0")
	s1ID := db.CreateStack("s1")
	s2ID := db.CreateStack("s2")

	expectedStatus := DatabaseStatus{
		ID:           "8cfa8cb55c92fa403369a13fd12a8e01",
		Name:         "db",
		NumberStacks: 3,
		Stacks:       []string{s0ID.String(), s2ID.String(), s1ID.String()},
	}

	if status := db.Status(); !reflect.DeepEqual(status, expectedStatus) {
		t.Errorf("status is %v, expected %v", status, expectedStatus)
	}
}

func TestDatabaseStatus_Empty(t *testing.T) {
	db := NewDatabase("db")

	expectedStatus := DatabaseStatus{
		ID:           "8cfa8cb55c92fa403369a13fd12a8e01",
		Name:         "db",
		NumberStacks: 0,
		Stacks:       []string{},
	}

	if status := db.Status(); !reflect.DeepEqual(status, expectedStatus) {
		t.Errorf("status is %#v, expected %#v", status, expectedStatus)
	}
}

func TestDatabaseStatusToJSON(t *testing.T) {
	databaseStatus := DatabaseStatus{
		ID:           "123456789",
		Name:         "db",
		NumberStacks: 3,
		Stacks:       []string{"stack1", "stack2", "stack3"},
	}

	expectedToJSON := `{"id":"123456789","name":"db","number_of_stacks":3,"stacks":["stack1","stack2","stack3"]}`

	if toJSON := databaseStatus.ToJSON(); string(toJSON) != expectedToJSON {
		t.Errorf("toJSON is %s, expected %s", string(toJSON), expectedToJSON)
	}

}

func TestDatabaseStatusToJSON_Empty(t *testing.T) {
	databaseStatus := DatabaseStatus{
		ID:           "123456789",
		Name:         "db",
		NumberStacks: 0,
	}

	expectedToJSON := `{"id":"123456789","name":"db","number_of_stacks":0}`

	if toJSON := databaseStatus.ToJSON(); string(toJSON) != expectedToJSON {
		t.Errorf("toJSON is %s, expected %s", string(toJSON), expectedToJSON)
	}

}

func TestDatabaseStacksStatus(t *testing.T) {
	s1 := NewStack("stack1")
	s1.Push("foo")

	s2 := NewStack("stack2")
	s2.Push(1)
	s2.Push(8)

	s3 := NewStack("stack3")

	db := NewDatabase("db")
	_ = db.AddStack(s1)
	_ = db.AddStack(s2)
	_ = db.AddStack(s3)

	expectedStatus := StacksStatus{
		Stacks: []StackStatus{s1.Status(), s2.Status(), s3.Status()},
	}

	if status := db.StacksStatus(); !reflect.DeepEqual(status, expectedStatus) {
		t.Errorf("status is %v, expected %v", status, expectedStatus)
	}
}

func TestDatabaseStacksStatus_Empty(t *testing.T) {
	db := NewDatabase("db")

	expectedStatus := StacksStatus{
		Stacks: []StackStatus{},
	}

	if status := db.StacksStatus(); !reflect.DeepEqual(status, expectedStatus) {
		t.Errorf("status is %v, expected %v", status, expectedStatus)
	}
}

func TestDatabaseStacksKV(t *testing.T) {
	s1 := NewStack("stack1")
	s1.Push("foo")

	s2 := NewStack("stack2")
	s2.Push(1)
	s2.Push(8)

	s3 := NewStack("stack3")

	db := NewDatabase("db")
	_ = db.AddStack(s1)
	_ = db.AddStack(s2)
	_ = db.AddStack(s3)

	expectedKV := StacksKV{
		Stacks: map[string]interface{}{
			"stack1": "foo",
			"stack2": 8,
			"stack3": nil,
		},
	}

	if kv := db.StacksKV(); !reflect.DeepEqual(kv, expectedKV) {
		t.Errorf("key-value is %v, expected %v", kv, expectedKV)
	}
}

func TestDatabaseStacksKV_Empty(t *testing.T) {
	db := NewDatabase("db")

	expectedKV := StacksKV{
		Stacks: map[string]interface{}{},
	}

	if kv := db.StacksKV(); !reflect.DeepEqual(kv, expectedKV) {
		t.Errorf("key-value is %v, expected %v", kv, expectedKV)
	}
}
