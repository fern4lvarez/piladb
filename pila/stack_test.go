package pila

import (
	"bytes"
	"fmt"
	"reflect"
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

func TestSetDatabase(t *testing.T) {
	db := NewDatabase("test-db")
	stack := NewStack("test-stack")
	stack.SetDatabase(db)

	if !reflect.DeepEqual(stack.Database, db) {
		t.Errorf("stack.Database is %v, expected %v", stack.Database, db)
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

func TestStackFlush(t *testing.T) {
	stack := NewStack("test-stack")
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

func TestStackSetID(t *testing.T) {
	db := NewDatabase("test-db")

	stack := NewStack("test-stack")
	stack.Database = db
	stack.SetID()

	if stack.ID.String() != "378c2601e338a49341d9858081452226" {
		t.Errorf("stack.ID is %s, expected %s", stack.ID.String(), "378c2601e338a49341d9858081452226")
	}
}

func TestStackSetID_NoDatabase(t *testing.T) {
	stack := NewStack("test-stack")
	stack.SetID()

	if stack.ID.String() != "2f44edeaa249ba81db20e9ddf000ba65" {
		t.Errorf("stack.ID is %s, expected %s", stack.ID.String(), "2f44edeaa249ba81db20e9ddf000ba65")
	}
}

func TestStackStatusJSON(t *testing.T) {
	stack := NewStack("test-stack")
	stack.Push("test")
	stack.Push(8)
	stack.Push(5.87)
	stack.Push([]byte("test"))

	expectedStatus := fmt.Sprintf(`{"id":"2f44edeaa249ba81db20e9ddf000ba65","name":"test-stack","peek":"dGVzdA==","size":4}`)
	if status, err := stack.Status().ToJSON(); err != nil {
		t.Fatal(err)
	} else if string(status) != expectedStatus {
		t.Errorf("status is %s, expected %s", string(status), expectedStatus)
	}
}

func TestStackStatusJSON_Empty(t *testing.T) {
	stack := NewStack("test-stack")

	expectedStatus := fmt.Sprintf(`{"id":"2f44edeaa249ba81db20e9ddf000ba65","name":"test-stack","peek":null,"size":0}`)
	if status, err := stack.Status().ToJSON(); err != nil {
		t.Fatal(err)
	} else if string(status) != expectedStatus {
		t.Errorf("status is %s, expected %s", string(status), expectedStatus)
	}
}

func TestStackStatusJSON_Error(t *testing.T) {
	// From https://golang.org/src/encoding/json/encode.go?s=5438:5481#L125
	// Channel, complex, and function values cannot be encoded in JSON.
	// Attempting to encode such a value causes Marshal to return
	// an UnsupportedTypeError.

	ch := make(chan int)
	f := func() string { return "a" }

	stack := NewStack("test-stack-channel")
	stack.Push(ch)

	if _, err := stack.Status().ToJSON(); err == nil {
		t.Error("err is nil, expected UnsupportedTypeError")
	}

	stack = NewStack("test-stack-function")
	stack.Push(f)

	if _, err := stack.Status().ToJSON(); err == nil {
		t.Error("err is nil, expected UnsupportedTypeError")
	}
}

func TestStacksStatusJSON(t *testing.T) {
	stack1 := NewStack("test-stack-1")
	stack1.Push("test")
	stack1.Push(8)
	stack1.Push(5.87)
	stack1.Push([]byte("test"))

	stack2 := NewStack("test-stack-2")
	stack2.Push("foo")
	stack2.Push([]byte("bar"))
	stack2.Push(999)

	stacksStatus := StacksStatus{
		Stacks: []StackStatus{stack1.Status(), stack2.Status()},
	}

	expectedStatus := fmt.Sprintf(`{"stacks":[{"id":"a0bfff209889f6f782997a7bd5b3d536","name":"test-stack-1","peek":"dGVzdA==","size":4},{"id":"f0d682fdfb3396c6f21e6f4d1d0da1cd","name":"test-stack-2","peek":999,"size":3}]}`)
	if status, err := stacksStatus.ToJSON(); err != nil {
		t.Fatal(err)
	} else if string(status) != expectedStatus {
		t.Errorf("status is %s, expected %s", string(status), expectedStatus)
	}
}

func TestStacksStatusJSON_Empty(t *testing.T) {
	stacksStatus := StacksStatus{
		Stacks: []StackStatus{},
	}

	expectedStatus := fmt.Sprintf(`{"stacks":[]}`)
	if status, err := stacksStatus.ToJSON(); err != nil {
		t.Fatal(err)
	} else if string(status) != expectedStatus {
		t.Errorf("status is %s, expected %s", string(status), expectedStatus)
	}
}

func TestStacksStatusJSON_Error(t *testing.T) {
	// From https://golang.org/src/encoding/json/encode.go?s=5438:5481#L125
	// Channel, complex, and function values cannot be encoded in JSON.
	// Attempting to encode such a value causes Marshal to return
	// an UnsupportedTypeError.

	ch := make(chan int)
	f := func() string { return "a" }

	stack := NewStack("test-stack-channel")
	stack.Push(ch)

	stacksStatus := StacksStatus{
		Stacks: []StackStatus{stack.Status()},
	}

	if _, err := stacksStatus.ToJSON(); err == nil {
		t.Error("err is nil, expected UnsupportedTypeError")
	}

	stack = NewStack("test-stack-function")
	stack.Push(f)

	stacksStatus = StacksStatus{
		Stacks: []StackStatus{stack.Status()},
	}

	if _, err := stacksStatus.ToJSON(); err == nil {
		t.Error("err is nil, expected UnsupportedTypeError")
	}
}

func TestStacksStatusLen(t *testing.T) {
	stack1 := NewStack("test-stack-1")
	stack2 := NewStack("test-stack-2")

	stacksStatus := StacksStatus{
		Stacks: []StackStatus{stack1.Status(), stack2.Status()},
	}

	expectedLen := 2
	if len := stacksStatus.Len(); len != expectedLen {
		t.Errorf("len is %d, expected %d", len, expectedLen)
	}
}

func TestStacksStatusLess(t *testing.T) {
	stack1 := NewStack("test-stack-1")
	stack2 := NewStack("test-stack-2")

	stacksStatus := StacksStatus{
		Stacks: []StackStatus{stack1.Status(), stack2.Status()},
	}

	if less := stacksStatus.Less(0, 1); !less {
		t.Errorf("less is %v, expected %v", less, true)
	}
}

func TestStacksStatusSwap(t *testing.T) {
	stack1 := NewStack("test-stack-1")
	stack2 := NewStack("test-stack-2")

	stacksStatus := StacksStatus{
		Stacks: []StackStatus{stack1.Status(), stack2.Status()},
	}

	expectedStacksStatus := StacksStatus{
		Stacks: []StackStatus{stack2.Status(), stack1.Status()},
	}

	if stacksStatus.Swap(0, 1); !reflect.DeepEqual(stacksStatus, expectedStacksStatus) {
		t.Errorf("status is %v, expected %v", stacksStatus, expectedStacksStatus)
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
