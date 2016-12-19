package pila

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/fern4lvarez/piladb/pkg/date"
)

func TestStackStatusJSON(t *testing.T) {
	now := time.Now().UTC()
	after := time.Now().UTC()

	stack := NewStack("test-stack", now)
	stack.Push("test")
	stack.Push(8)
	stack.Push(5.87)
	stack.Push([]byte("test"))
	stack.Update(after)

	expectedStatus := fmt.Sprintf(`{"id":"2f44edeaa249ba81db20e9ddf000ba65","name":"test-stack","peek":"dGVzdA==","size":4,"created_at":"%v","updated_at":"%v","read_at":"%v"}`,
		date.Format(now.Local()),
		date.Format(after.Local()),
		date.Format(after.Local()))
	if status, err := stack.Status().ToJSON(); err != nil {
		t.Fatal(err)
	} else if string(status) != expectedStatus {
		t.Errorf("status is %s, expected %s", string(status), expectedStatus)
	}
}

func TestStackStatusJSON_Empty(t *testing.T) {
	now := time.Now().UTC()
	stack := NewStack("test-stack", now)
	stack.Update(now)

	expectedStatus := fmt.Sprintf(`{"id":"2f44edeaa249ba81db20e9ddf000ba65","name":"test-stack","peek":null,"size":0,"created_at":"%v","updated_at":"%v","read_at":"%v"}`,
		date.Format(now.Local()),
		date.Format(now.Local()),
		date.Format(now.Local()))
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

	stack := NewStack("test-stack-channel", time.Now())
	stack.Push(ch)

	if _, err := stack.Status().ToJSON(); err == nil {
		t.Error("err is nil, expected UnsupportedTypeError")
	}

	stack = NewStack("test-stack-function", time.Now())
	stack.Push(f)

	if _, err := stack.Status().ToJSON(); err == nil {
		t.Error("err is nil, expected UnsupportedTypeError")
	}
}

func TestStacksStatusJSON(t *testing.T) {
	now := time.Now().UTC().UTC()
	after := time.Now().UTC().UTC()

	stack1 := NewStack("test-stack-1", now)
	stack1.Push("test")
	stack1.Push(8)
	stack1.Push(5.87)
	stack1.Push([]byte("test"))
	stack1.Update(after)

	stack2 := NewStack("test-stack-2", now)
	stack2.Push("foo")
	stack2.Push([]byte("bar"))
	stack2.Push(999)
	stack2.Update(after)

	stacksStatus := StacksStatus{
		Stacks: []StackStatus{stack1.Status(), stack2.Status()},
	}

	expectedStatus := fmt.Sprintf(`{"stacks":[{"id":"a0bfff209889f6f782997a7bd5b3d536","name":"test-stack-1","peek":"dGVzdA==","size":4,"created_at":"%v","updated_at":"%v","read_at":"%v"},{"id":"f0d682fdfb3396c6f21e6f4d1d0da1cd","name":"test-stack-2","peek":999,"size":3,"created_at":"%v","updated_at":"%v","read_at":"%v"}]}`,
		date.Format(now.Local()), date.Format(after.Local()), date.Format(after.Local()),
		date.Format(now.Local()), date.Format(after.Local()), date.Format(after.Local()))
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

	input := []interface{}{ch, f}

	for _, in := range input {
		stack := NewStack("test-stack", time.Now().UTC())
		stack.Push(in)

		stacksStatus := StacksStatus{
			Stacks: []StackStatus{stack.Status()},
		}

		if _, err := stacksStatus.ToJSON(); err == nil {
			t.Error("err is nil, expected UnsupportedTypeError")
		}
	}
}

func TestStacksStatusLen(t *testing.T) {
	stack1 := NewStack("test-stack-1", time.Now().UTC())
	stack2 := NewStack("test-stack-2", time.Now().UTC())

	stacksStatus := StacksStatus{
		Stacks: []StackStatus{stack1.Status(), stack2.Status()},
	}

	expectedLen := 2
	if len := stacksStatus.Len(); len != expectedLen {
		t.Errorf("len is %d, expected %d", len, expectedLen)
	}
}

func TestStacksStatusLess(t *testing.T) {
	stack1 := NewStack("test-stack-1", time.Now().UTC())
	stack2 := NewStack("test-stack-2", time.Now().UTC())

	stacksStatus := StacksStatus{
		Stacks: []StackStatus{stack1.Status(), stack2.Status()},
	}

	if less := stacksStatus.Less(0, 1); !less {
		t.Errorf("less is %v, expected %v", less, true)
	}
}

func TestStacksStatusSwap(t *testing.T) {
	stack1 := NewStack("test-stack-1", time.Now().UTC())
	stack2 := NewStack("test-stack-2", time.Now().UTC())

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

func TestStacksKVJSON(t *testing.T) {
	stack1 := NewStack("test-stack-1", time.Now().UTC())
	stack1.Push("test")

	stack2 := NewStack("test-stack-2", time.Now().UTC())
	stack2.Push(999)

	stacksKV := StacksKV{
		Stacks: map[string]interface{}{
			"test-stack-1": "test",
			"test-stack-2": 999,
		},
	}

	expectedKV := fmt.Sprintf(`{"stacks":{"test-stack-1":"test","test-stack-2":999}}`)
	if skv, err := stacksKV.ToJSON(); err != nil {
		t.Fatal(err)
	} else if string(skv) != expectedKV {
		t.Errorf("stacks key-value is %s, expected %s", string(skv), expectedKV)
	}
}

func TestStacksKVJSON_Error(t *testing.T) {
	// From https://golang.org/src/encoding/json/encode.go?s=5438:5481#L125
	// Channel, complex, and function values cannot be encoded in JSON.
	// Attempting to encode such a value causes Marshal to return
	// an UnsupportedTypeError.
	ch := make(chan int)
	f := func() string { return "a" }

	input := []interface{}{ch, f}

	for _, in := range input {
		stack := NewStack("test-stack", time.Now().UTC())
		stack.Push(in)

		stacksKV := StacksKV{
			Stacks: map[string]interface{}{
				stack.Name: stack.Peek(),
			},
		}

		if _, err := stacksKV.ToJSON(); err == nil {
			t.Error("err is nil, expected UnsupportedTypeError")
		}
	}
}
