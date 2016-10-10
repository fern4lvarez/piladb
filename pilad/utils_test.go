package main

import (
	"reflect"
	"testing"

	"github.com/fern4lvarez/piladb/pila"
	"github.com/fern4lvarez/piladb/pkg/uuid"
)

func TestResourceDatabase(t *testing.T) {
	dbName := "db"
	inputs := []string{dbName, uuid.New(dbName).String()}

	for _, input := range inputs {
		expectedDB := pila.NewDatabase(dbName)

		p := pila.NewPila()
		_ = p.AddDatabase(expectedDB)

		conn := NewConn()
		conn.Pila = p

		db, ok := ResourceDatabase(conn, input)
		if !ok {
			t.Errorf("ok is %v, expected true", ok)
		}
		if !reflect.DeepEqual(expectedDB, db) {
			t.Errorf("db is %v, expected db id %v", db, expectedDB)
		}
	}
}

func TestResourceDatabase_False(t *testing.T) {
	dbName := "db"
	inputs := []string{dbName, uuid.New(dbName).String()}

	for _, input := range inputs {
		p := pila.NewPila()

		conn := NewConn()
		conn.Pila = p

		_, ok := ResourceDatabase(conn, input)
		if ok {
			t.Errorf("ok is %v, expected false", ok)
		}
	}
}

func TestResourceStack(t *testing.T) {
	dbName := "db"

	stackName := "stack"
	inputs := []string{stackName, uuid.New(dbName + stackName).String()}

	for _, input := range inputs {
		expectedStack := pila.NewStack(stackName)

		db := pila.NewDatabase(dbName)
		_ = db.AddStack(expectedStack)

		stack, ok := ResourceStack(db, input)
		if !ok {
			t.Errorf("ok is %v, expected true", ok)
		}
		if !reflect.DeepEqual(expectedStack, stack) {
			t.Errorf("stack is %v, expected stack id %v", stack, expectedStack)
		}
	}
}

func TestResourceStack_False(t *testing.T) {
	dbName := "db"

	stackName := "stack"
	inputs := []string{stackName, uuid.New(dbName + stackName).String()}

	for _, input := range inputs {
		db := pila.NewDatabase(dbName)

		_, ok := ResourceStack(db, input)
		if ok {
			t.Errorf("ok is %v, expected false", ok)
		}
	}
}

func TestMemStats(t *testing.T) {
	memStats := MemStats()

	if memStats == nil {
		t.Fatalf("memory stats is nil")
	}

	if memStats.Alloc < 0 {
		t.Errorf("memory allocated is negative, expected positive")
	}
}
