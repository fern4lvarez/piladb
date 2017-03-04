package pila

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

type TestBaseStack struct{}

func (s *TestBaseStack) Push(element interface{}) { return }
func (s *TestBaseStack) Pop() (interface{}, bool) { return nil, false }
func (s *TestBaseStack) Size() int                { return 0 }
func (s *TestBaseStack) Peek() interface{}        { return nil }
func (s *TestBaseStack) Flush()                   { return }
func (s *TestBaseStack) Block() bool              { return false }

func TestNewStack(t *testing.T) {
	now := time.Now()
	stack := NewStack("test-stack", now)

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
	if stack.CreatedAt != now {
		t.Errorf("stack.CreatedAt is %v, expected %v", stack.CreatedAt, now)
	}
	if stack.base == nil {
		t.Fatalf("stack.base is nil")
	}
	if stack.base.Size() != 0 {
		t.Fatalf("stack.base.Size() is %d, expected %d", stack.base.Size(), 0)
	}
}

func TestNewStackWithBase(t *testing.T) {
	now := time.Now()
	stack := NewStackWithBase("test-stack", now, &TestBaseStack{})

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
	if stack.CreatedAt != now {
		t.Errorf("stack.CreatedAt is %v, expected %v", stack.CreatedAt, now)
	}
	if stack.base == nil {
		t.Fatalf("stack.base is nil")
	}
	if stack.base.Size() != 0 {
		t.Fatalf("stack.base.Size() is %d, expected %d", stack.base.Size(), 0)
	}
}

func TestSetDatabase(t *testing.T) {
	db := NewDatabase("test-db")
	stack := NewStack("test-stack", time.Now())
	stack.SetDatabase(db)

	if !reflect.DeepEqual(stack.Database, db) {
		t.Errorf("stack.Database is %v, expected %v", stack.Database, db)
	}
}

func TestStackPush(t *testing.T) {
	stack := NewStack("test-stack", time.Now())
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
	stack := NewStack("test-stack", time.Now())
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
	stack := NewStack("test-stack", time.Now())
	_, ok := stack.Pop()
	if ok {
		t.Error("stack.Pop() is ok")
	}
}

func TestStackSize(t *testing.T) {
	stack := NewStack("test-stack", time.Now())
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

func TestStackBlock(t *testing.T) {
	stack := NewStack("test-stack", time.Now())
	stack.Push("test")
	stack.Push(8)
	stack.Block()

	if stack.Blocked != true {
		t.Errorf("stack isn't blocked")
	}
}

func TestStackPeek(t *testing.T) {
	stack := NewStack("test-stack", time.Now())
	stack.Push("test")
	stack.Push(8)

	element := stack.Peek()
	if element != 8 {
		t.Errorf("element is %v, expected %v", element, 8)
	}
}

func TestStackFlush(t *testing.T) {
	stack := NewStack("test-stack", time.Now())
	stack.Push("test")
	stack.Push(8)
	stack.Push(87.443)

	stack.Flush()
	if stack.Size() != 0 {
		t.Errorf("stack is not empty")
	}
	if stack.Peek() != nil {
		t.Errorf("stack peek is not nil")
	}
}

func TestStackUpdate(t *testing.T) {
	now := time.Now()
	updateTime := time.Now()
	stack := NewStack("test-stack", now)

	stack.Update(updateTime)

	if stack.CreatedAt != now {
		t.Errorf("stack.CreatedAt is %v, expected %v", stack.CreatedAt, now)
	}
	if stack.UpdatedAt != updateTime {
		t.Errorf("stack.UpdatedAt is %v, expected %v", stack.UpdatedAt, updateTime)
	}
	if stack.ReadAt != updateTime {
		t.Errorf("stack.ReadAt is %v, expected %v", stack.UpdatedAt, updateTime)
	}
}

func TestStackRead(t *testing.T) {
	now := time.Now()
	updateTime := time.Now()
	readTime := time.Now()
	stack := NewStack("test-stack", now)

	stack.Update(updateTime)
	stack.Read(readTime)

	if stack.CreatedAt != now {
		t.Errorf("stack.CreatedAt is %v, expected %v", stack.CreatedAt, now)
	}
	if stack.UpdatedAt != updateTime {
		t.Errorf("stack.UpdatedAt is %v, expected %v", stack.UpdatedAt, updateTime)
	}
	if stack.ReadAt != readTime {
		t.Errorf("stack.ReadAt is %v, expected %v", stack.UpdatedAt, updateTime)
	}
}

func TestStackSetID(t *testing.T) {
	db := NewDatabase("test-db")

	stack := NewStack("test-stack", time.Now())
	stack.Database = db
	stack.SetID()

	if stack.ID.String() != "378c2601e338a49341d9858081452226" {
		t.Errorf("stack.ID is %s, expected %s", stack.ID.String(), "378c2601e338a49341d9858081452226")
	}
}

func TestStackSetID_NoDatabase(t *testing.T) {
	stack := NewStack("test-stack", time.Now())
	stack.SetID()

	if stack.ID.String() != "2f44edeaa249ba81db20e9ddf000ba65" {
		t.Errorf("stack.ID is %s, expected %s", stack.ID.String(), "2f44edeaa249ba81db20e9ddf000ba65")
	}
}

func TestStackSizeToJSON(t *testing.T) {
	stack := NewStack("test-stack", time.Now())
	stack.Push("test")
	stack.Push(8)
	stack.Push(87.443)

	expectedSize := `3`
	if string(stack.SizeToJSON()) != expectedSize {
		t.Errorf("size is %s, expected %s", string(stack.SizeToJSON()), expectedSize)
	}
}

func TestElementJSON(t *testing.T) {
	elements := []Element{
		{Value: "foo"},
		{Value: 42},
		{Value: 3.14},
		{Value: []byte("hello")},
		{Value: map[string]int{"one": 1}},
	}
	expectedElements := []string{
		`{"element":"foo"}`,
		`{"element":42}`,
		`{"element":3.14}`,
		`{"element":"aGVsbG8="}`,
		`{"element":{"one":1}}`,
	}

	for i, element := range elements {
		expectedElement := expectedElements[i]
		if element, err := element.ToJSON(); err != nil {
			t.Fatal(err)
		} else if string(element) != expectedElement {
			t.Errorf("element is %s, expected %s", string(element), expectedElement)
		}
	}

}

func TestElementJSON_Error(t *testing.T) {
	// From https://golang.org/src/encoding/json/encode.go?s=5438:5481#L125
	// Channel, complex, and function values cannot be encoded in JSON.
	// Attempting to encode such a value causes Marshal to return
	// an UnsupportedTypeError.

	ch := make(chan int)
	f := func() string { return "a" }
	elements := []Element{
		{Value: ch},
		{Value: f},
	}

	for _, element := range elements {
		if _, err := element.ToJSON(); err == nil {
			t.Error("err is nil, expected UnsupportedTypeError")
		}
	}
}

func TestElementDecode(t *testing.T) {
	elementReaders := []string{
		`{"element":"foo"}`,
		`{"element":42}`,
		`{"element":3.14}`,
		`{"element":"aGVsbG8="}`,
		`{"element":{"one":1}}`,
	}
	expectedElements := []Element{
		{Value: "foo"},
		{Value: 42.0}, // decode int into float
		{Value: 3.14},
		{Value: "aGVsbG8="},                         // does not decode into []byte
		{Value: map[string]interface{}{"one": 1.0}}, // decode inner int into float
	}

	for i, elementReader := range elementReaders {
		expectedElement := expectedElements[i]
		r := bytes.NewBuffer([]byte(elementReader))

		var element Element
		if err := element.Decode(r); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(element.Value, expectedElement.Value) {
			t.Errorf("element is %#v, expected %#v", element.Value, expectedElement.Value)
		}
	}
}

func TestElementDecode_Error(t *testing.T) {
	elementReaders := []string{
		`{`,
		`}`,
		``,
		` `,
		`$`,
		`%{}`,
	}

	for _, elementReader := range elementReaders {
		r := bytes.NewBuffer([]byte(elementReader))

		var element Element
		if err := element.Decode(r); err == nil {
			t.Fatal("err is nil, expected error")
		}
	}
}
