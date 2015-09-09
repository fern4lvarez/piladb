package pila

import (
	"reflect"
	"testing"
)

func TestNewPila(t *testing.T) {
	pila := NewPila()
	if pila == nil {
		t.Error("pila is nil")
	} else if pila.Databases == nil {
		t.Error("pila.Databases is nil")
	}

}

func TestPilaCreateDatabase(t *testing.T) {
	pila := NewPila()
	id := pila.CreateDatabase("test-1")

	db, ok := pila.Databases[id]
	if !ok {
		t.Errorf("db %v not added to pila", id)
	} else if !reflect.DeepEqual(db.Pila, pila) {
		t.Errorf("db %v does not contain pile", id)
	}
}
