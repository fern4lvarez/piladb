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

func TestPilaAddDatabase(t *testing.T) {
	pila := NewPila()
	db := NewDatabase("test")
	err := pila.AddDatabase(db)
	if err != nil {
		t.Error("err is not nil")
	}

	id := db.ID
	db, ok := pila.Databases[id]
	if !ok {
		t.Errorf("db %v not added to pila", id)
	} else if !reflect.DeepEqual(db.Pila, pila) {
		t.Errorf("db %v does not contain pile", id)
	}

}

func TestPilaAddDatabase_Error(t *testing.T) {
	pila := NewPila()
	pila2 := NewPila()
	db := NewDatabase("test")
	db2 := NewDatabase("test")

	if err := pila.AddDatabase(db); err != nil {
		t.Errorf("err is not nil, but %v", err)
	}

	if err := pila2.AddDatabase(db); err == nil {
		t.Error("err is nil")
	}

	if err := pila.AddDatabase(db2); err == nil {
		t.Error("err is nil")
	}
}
