package pila

import (
	"encoding/json"
	"time"
)

// StackStatuser represents an interface for
// types representing the status of stacks or stacks
// in different formats.
type StackStatuser interface {
	ToJSON() ([]byte, error)
}

// StackStatus represents the status of a Stack.
type StackStatus struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Peek      interface{} `json:"peek"`
	Size      int         `json:"size"`
	Blocked   bool        `json:"blocked"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	ReadAt    time.Time   `json:"read_at"`
}

// ToJSON converts a StackStatus into JSON.
func (stackStatus StackStatus) ToJSON() ([]byte, error) {
	return json.Marshal(stackStatus)
}

// StacksStatus represents the status of a list of Stacks.
type StacksStatus struct {
	Stacks []StackStatus `json:"stacks"`
}

// ToJSON converts a StacksStatus into JSON.
func (stacksStatus StacksStatus) ToJSON() ([]byte, error) {
	return json.Marshal(stacksStatus)
}

// Len return the length of the list of Stacks.
func (stacksStatus StacksStatus) Len() int {
	return len(stacksStatus.Stacks)
}

// Less determines whether a StackStatus on the list is less than other.
func (stacksStatus StacksStatus) Less(i, j int) bool {
	return stacksStatus.Stacks[i].Name < stacksStatus.Stacks[j].Name
}

// Swap swaps positions between two StackStatus.
func (stacksStatus StacksStatus) Swap(i, j int) {
	stacksStatus.Stacks[i], stacksStatus.Stacks[j] = stacksStatus.Stacks[j], stacksStatus.Stacks[i]
}

// StacksKV represents a list of status by a key-value list
// composed by name and peek of the Stack.
type StacksKV struct {
	Stacks map[string]interface{} `json:"stacks"`
}

// ToJSON converts a StacksKV into JSON.
func (stacksKV StacksKV) ToJSON() ([]byte, error) {
	return json.Marshal(stacksKV)
}
