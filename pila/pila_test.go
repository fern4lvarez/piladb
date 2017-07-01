package pila

import (
	"reflect"
	"testing"
)

func TestNewPila(t *testing.T) {
	pila := NewPila()
	if pila == nil {
		t.Fatal("pila is nil")
	}
	if pila.Databases == nil {
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
		t.Fatal("err is not nil")
	}

	id := db.ID
	db, ok := pila.Databases[id]
	if !ok {
		t.Errorf("db %v not added to pila", id)
	} else if !reflect.DeepEqual(db.Pila, pila) {
		t.Errorf("db %v does not contain pila", id)
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

func TestPilaRemoveDatabase(t *testing.T) {
	pila := NewPila()
	db := NewDatabase("test")
	pila.AddDatabase(db)

	if ok := pila.RemoveDatabase(db.ID); !ok {
		t.Errorf("RemoveDatabase did not succeed")
	}

	if db.Pila != nil {
		t.Errorf("a pila is assigned to database %v", db.Name)
	}

	if _, ok := pila.Databases[db.ID]; ok {
		t.Errorf("Removed database does exist on pila")
	}
}

func TestPilaRemoveDatabase_False(t *testing.T) {
	pila := NewPila()
	db := NewDatabase("test")

	if ok := pila.RemoveDatabase(db.ID); ok {
		t.Errorf("RemoveDatabase did succeed")
	}
}

func TestPilaDatabase(t *testing.T) {
	pila := NewPila()
	db := NewDatabase("test")
	err := pila.AddDatabase(db)
	if err != nil {
		t.Fatal("err is not nil")
	}

	db2, ok := pila.Database(db.ID)
	if !ok {
		t.Errorf("pila has no Database %v", db.ID)
	} else if !reflect.DeepEqual(db2, db) {
		t.Errorf("Database is %v, expected %v", db2, db)
	}
}

func TestPilaDatabase_False(t *testing.T) {
	pila := NewPila()
	db := NewDatabase("test")
	db2, ok := pila.Database(db.ID)
	if ok {
		t.Errorf("pila has Database %v", db.ID)
	} else if db2 != nil {
		t.Errorf("Database %v is not nil", db2)
	}
}

func TestPilaStatusToJSON(t *testing.T) {
	pila := NewPila()
	db0 := NewDatabase("db0")
	pila.AddDatabase(db0)

	expectedStatus := `{"number_of_databases":1,"databases":[{"id":"91010edc-36f6-25cc-5b10-f2648eb2b322","name":"db0","number_of_stacks":0}]}`

	if status := pila.Status().ToJSON(); string(status) != expectedStatus {
		t.Errorf("status is %s, expected %s", string(status), expectedStatus)
	}
}

func TestPilaStatusToJSON_Empty(t *testing.T) {
	pila := NewPila()

	expectedStatus := `{"number_of_databases":0,"databases":[]}`

	if status := pila.Status().ToJSON(); string(status) != expectedStatus {
		t.Errorf("status is %s, expected %s", string(status), expectedStatus)
	}
}
