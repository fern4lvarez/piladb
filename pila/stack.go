package pila

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/fern4lvarez/piladb/pkg/stack"
	"github.com/fern4lvarez/piladb/pkg/uuid"
)

// Stack represents a stack entity in piladb.
type Stack struct {
	// ID is a unique identifier of the Stack
	// Note: Do not use this field to read the ID,
	// as this method is not thread-safe. See UUID()
	// instead.
	ID fmt.Stringer

	// Name of the Stack
	Name string

	// Database associated to the Stack
	Database *Database

	// CreatedAt represents the date when the Stack was created
	CreatedAt time.Time

	// UpdatedAt represents the date when the Stack was updated for the last time.
	// This date must be updated when a Stack is created, and when receives a PUSH,
	// POP, or FLUSH operation.
	// Note that unlike CreatedAt, UpdatedAt is not triggered automatically
	// when one of these events happens, but it needs to be set by hand.
	UpdatedAt time.Time

	// ReadAt represents the date when the Stack was read for the last time.
	// This date must be updated when a Stack is created, accessed, and when it
	// receives a PUSH, POP, or FLUSH operation.
	// Note that unlike CreatedAt, ReadAt is not triggered automatically
	// when one of these events happens, but it needs to be set by hand.
	ReadAt time.Time

	// blocked specifies whether a Stack is blocked, so it can be only read, but
	// not modified.
	blocked bool

	// dateMu serves as a mutex to lock dates on concurrent
	// updates in order to avoid race conditions.
	dateMu sync.Mutex

	// IDMu provides a mutex to handle concurrent reads and
	// writes on the Stack ID.
	IDMu sync.RWMutex

	// blockMu provides a mutex to handle concurrent reads and
	// writes on the Stack Blocked status.
	blockMu sync.RWMutex

	// base represents the Stack data structure
	base stack.Stacker
}

// NewStack creates a new Stack given a name and a creation date,
// without an association to any Database. It uses the default
// ./pkg/stack implementation as a base Stack.
func NewStack(name string, t time.Time) *Stack {
	return NewStackWithBase(name, t, stack.NewStack())
}

// NewStackWithBase creates a new Stack given a name, a creation date,
// and a stack.Stacker base implementation, without an association to any Database.
func NewStackWithBase(name string, t time.Time, base stack.Stacker) *Stack {
	s := &Stack{}
	s.Name = name
	s.SetID()
	s.CreatedAt = t
	s.base = base
	return s
}

// Push an element on top of the Stack.
func (s *Stack) Push(element interface{}) error {
	if s.Blocked() {
		return errors.New("Stack is blocked")
	}
	s.base.Push(element)
	return nil
}

// Pop removes and returns the element on top of the Stack.
// If the Stack was empty, it returns false.
func (s *Stack) Pop() (interface{}, error) {
	if s.Blocked() {
		return nil, errors.New("Stack is blocked")
	}

	element := s.base.Pop()

	if element != nil {
		return element, nil
	}
	return nil, errors.New("Stack is empty")
}

// Base bases the Stack on an element, so this becomes
// the bottommost one of the Stack.
func (s *Stack) Base(element interface{}) error {
	if s.Blocked() {
		return errors.New("Stack is blocked")
	}

	s.base.Base(element)
	return nil
}

// Sweep removes and returns the bottommost element of the Stack.
// If the Stack is empty or blocked, it returns an error.
func (s *Stack) Sweep() (interface{}, error) {
	if s.Blocked() {
		return nil, errors.New("Stack is blocked")
	}

	element, ok := s.base.Sweep()
	if !ok {
		return nil, errors.New("Stack is empty")
	}

	return element, nil
}

// SweepPush removes and returns the bottommost element of the Stack,
// and pushes an element on top of it, as an atomic operation.
// If the Stack is empty or blocked, it returns an error.
func (s *Stack) SweepPush(element interface{}) (interface{}, error) {
	if s.Blocked() {
		return nil, errors.New("Stack is blocked")
	}

	element, ok := s.base.SweepPush(element)
	if !ok {
		return nil, errors.New("Stack is empty")
	}

	return element, nil
}

// Rotate moves the bottommost element of the Stack
// to the top. If the Stack is empty or blocked,
// it returns an error.
func (s *Stack) Rotate() error {
	if s.Blocked() {
		return errors.New("Stack is blocked")
	}

	if !s.base.Rotate() {
		return errors.New("Stack is empty")
	}

	return nil
}

// Size returns the size of the Stack.
func (s *Stack) Size() int {
	return s.base.Size()
}

// Empty returns true if a stack is empty.
func (s *Stack) Empty() bool {
	return s.base.Size() == 0
}

// Peek returns the element on top of the Stack.
func (s *Stack) Peek() interface{} {
	return s.base.Peek()
}

// Flush flushes the content of the Stack.
func (s *Stack) Flush() error {
	if s.Blocked() {
		return errors.New("Stack is blocked")
	}
	s.base.Flush()
	return nil
}

// Update takes a date and updates UpdateAt and ReadAt
// fields of the Stack.
func (s *Stack) Update(t time.Time) {
	s.dateMu.Lock()
	s.UpdatedAt = t
	s.ReadAt = t
	s.dateMu.Unlock()
}

// Read takes a date and updates ReadAt field
// of the Stack.
func (s *Stack) Read(t time.Time) {
	s.dateMu.Lock()
	s.ReadAt = t
	s.dateMu.Unlock()
}

// SetDatabase links the Stack with a given Database and
// recalculates its ID.
func (s *Stack) SetDatabase(db *Database) {
	s.Database = db
	s.SetID()
}

// UUID returns the unique Stack ID providing thread safety.
func (s *Stack) UUID() fmt.Stringer {
	s.IDMu.RLock()
	defer s.IDMu.RUnlock()

	return s.ID
}

// SetID recalculates the id of the Stack based on its
// Database name and its own name.
func (s *Stack) SetID() {
	s.IDMu.Lock()
	defer s.IDMu.Unlock()

	if s.Database != nil {
		s.ID = uuid.New(s.Database.Name + s.Name)
		return
	}

	s.ID = uuid.New(s.Name)
}

// SizeToJSON returns the size of the Stack encoded as json.
func (s *Stack) SizeToJSON() []byte {
	// Do not check error as we consider the size
	// of a stack valid for a JSON encoding.
	size, _ := json.Marshal(s.Size())
	return size
}

// Status returns the status of the Stack  in json format.
func (s *Stack) Status() StackStatus {
	status := StackStatus{}
	status.ID = s.UUID().String()
	status.Name = s.Name
	status.Size = s.Size()
	status.Peek = s.Peek()
	status.Blocked = s.Blocked()
	status.CreatedAt = s.CreatedAt.Local()
	status.UpdatedAt = s.UpdatedAt.Local()
	status.ReadAt = s.ReadAt.Local()

	return status
}

// Block blocks the Stack.
func (s *Stack) Block() {
	s.blockMu.Lock()
	s.blocked = true
	s.blockMu.Unlock()
}

// Unblock unblocks the Stack.
func (s *Stack) Unblock() {
	s.blockMu.Lock()
	s.blocked = false
	s.blockMu.Unlock()
}

// Blocked returns true if Stack is blocked.
func (s *Stack) Blocked() bool {
	s.blockMu.RLock()
	defer s.blockMu.RUnlock()

	return s.blocked
}

// Element represents the payload of a Stack element.
type Element struct {
	Value interface{} `json:"element"`
}

// ToJSON converts an Element into JSON.
func (element Element) ToJSON() ([]byte, error) {
	return json.Marshal(element)
}

// Decode decodes json data into an Element.
// If the payload is nil, empty, or doesn't
// follow a {"element":...} pattern, it will
// return an error.
func (element *Element) Decode(r io.Reader) error {
	if r == nil {
		return errors.New("payload is nil")
	}

	elementBuffer := new(bytes.Buffer)
	elementBuffer.ReadFrom(r)

	if len(elementBuffer.Bytes()) == 0 {
		return errors.New("payload is empty")
	}

	if !bytes.HasPrefix(elementBuffer.Bytes(), []byte(`{"element"`)) {
		return errors.New("payload is malformed, missing element key?")
	}

	decoder := json.NewDecoder(elementBuffer)
	return decoder.Decode(element)
}
