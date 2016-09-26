package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fern4lvarez/piladb/pila"
	"github.com/fern4lvarez/piladb/pkg/uuid"
)

func TestNewConn(t *testing.T) {
	conn := NewConn()

	if conn == nil {
		t.Fatal("conn is nil")
	}

	if conn.Pila == nil {
		t.Error("conn.Pila is nil")
	}

	if conn.Status == nil {
		t.Fatal("conn.Status is nil")
	}

	if conn.Status.Code != "OK" {
		t.Errorf("conn.Status is %s, expected %s", conn.Status.Code, "OK")
	}
}

func TestStatusHandler(t *testing.T) {
	conn := NewConn()
	request, err := http.NewRequest("GET", "/_status", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	conn.statusHandler(response, request)

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
	}

	if response.Code != http.StatusOK {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusOK)
	}

	statusJSON, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(string(statusJSON), `{"status":"OK","version"`) {
		t.Errorf("status is %s", string(statusJSON))
	}
}

func TestDatabasesHandler_GET(t *testing.T) {
	db := pila.NewDatabase("db")

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	request, err := http.NewRequest("GET", "/databases", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	conn.databasesHandler(response, request)

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
	}

	if response.Code != http.StatusOK {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusOK)
	}

	databases, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if expected := `{"number_of_databases":1,"databases":[{"id":"8cfa8cb55c92fa403369a13fd12a8e01","name":"db","number_of_stacks":0}]}`; string(databases) != expected {
		t.Errorf("databases are %s, expected %s", string(databases), expected)
	}
}

func TestDatabasesHandler_GET_Empty(t *testing.T) {
	conn := NewConn()
	request, err := http.NewRequest("GET", "/databases", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	conn.databasesHandler(response, request)

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
	}

	if response.Code != http.StatusOK {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusOK)
	}

	databases, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(databases) != `{"number_of_databases":0,"databases":[]}` {
		t.Errorf("databases are %s, expected %s", string(databases), `{"number_of_databases":0,"databases":[]}`)
	}
}

func TestDatabasesHandler_PUT(t *testing.T) {
	// see TestCreateDatabaseHandler* tests for a further spec on
	// how createDatabaseHandler works.
	conn := NewConn()
	request, err := http.NewRequest("PUT", "/databases", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	conn.databasesHandler(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusBadRequest)
	}
}

func TestCreateDatabaseHandler(t *testing.T) {
	conn := NewConn()
	request, err := http.NewRequest("PUT", "/databases?name=db", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	conn.createDatabaseHandler(response, request)

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
	}

	if response.Code != http.StatusCreated {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusCreated)
	}

	databases, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(databases) != `{"id":"8cfa8cb55c92fa403369a13fd12a8e01","name":"db","number_of_stacks":0}` {
		t.Errorf("databases are %s, expected %s", string(databases), `{"id":"8cfa8cb55c92fa403369a13fd12a8e01","name":"db","number_of_stacks":0}`)
	}
}

func TestCreateDatabaseHandler_NoName(t *testing.T) {
	conn := NewConn()
	request, err := http.NewRequest("PUT", "/databases", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	conn.createDatabaseHandler(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusBadRequest)
	}
}

func TestCreateDatabaseHandler_Duplicated(t *testing.T) {
	conn := NewConn()

	request, err := http.NewRequest("PUT", "/databases?name=db", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	conn.createDatabaseHandler(response, request)

	if response.Code != http.StatusCreated {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusCreated)
	}

	request, err = http.NewRequest("PUT", "/databases?name=db", nil)
	if err != nil {
		t.Fatal(err)
	}
	response = httptest.NewRecorder()

	conn.createDatabaseHandler(response, request)

	if response.Code != http.StatusConflict {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusConflict)
	}
}

func TestDatabaseHandler_GET(t *testing.T) {
	s := pila.NewStack("stack")
	s.Push("foo")

	db := pila.NewDatabase("mydb")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	request, err := http.NewRequest("GET",
		fmt.Sprintf("/databases/%s",
			db.ID.String()),
		nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()

	databaseHandle := conn.databaseHandler(db.ID.String())
	databaseHandle.ServeHTTP(response, request)

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
	}

	if response.Code != http.StatusOK {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusOK)
	}

	database, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if expected := `{"id":"c13cec0e70876381c78c616ee2d809eb","name":"mydb","number_of_stacks":1,"stacks":["b92f53fa3884305ef798fd8c5d7609ad"]}`; string(database) != expected {
		t.Errorf("database is %v, expected %v", string(database), expected)
	}
}

func TestDatabaseHandler_GET_Name(t *testing.T) {
	s := pila.NewStack("stack")
	s.Push("foo")

	db := pila.NewDatabase("mydb")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	request, err := http.NewRequest("GET",
		fmt.Sprintf("/databases/%s",
			db.Name),
		nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()

	databaseHandle := conn.databaseHandler(db.Name)
	databaseHandle.ServeHTTP(response, request)

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
	}

	if response.Code != http.StatusOK {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusOK)
	}

	database, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if expected := `{"id":"c13cec0e70876381c78c616ee2d809eb","name":"mydb","number_of_stacks":1,"stacks":["b92f53fa3884305ef798fd8c5d7609ad"]}`; string(database) != expected {
		t.Errorf("database is %v, expected %v", string(database), expected)
	}
}

func TestDatabaseHandler_DELETE(t *testing.T) {
	db1 := pila.NewDatabase("mydb1")
	db2 := pila.NewDatabase("mydb2")

	p := pila.NewPila()
	_ = p.AddDatabase(db1)
	_ = p.AddDatabase(db2)

	conn := NewConn()
	conn.Pila = p

	request, err := http.NewRequest("DELETE",
		fmt.Sprintf("/databases/%s",
			db1.ID.String()),
		nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()

	databaseHandle := conn.databaseHandler(db1.ID.String())
	databaseHandle.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusNoContent)
	}

	if len(conn.Pila.Databases) != 1 {
		t.Errorf("got %d database, expected %d", len(conn.Pila.Databases), 1)
	}
}

func TestDatabaseHandler_DELETE_Name(t *testing.T) {
	db1 := pila.NewDatabase("mydb1")
	db2 := pila.NewDatabase("mydb2")

	p := pila.NewPila()
	_ = p.AddDatabase(db1)
	_ = p.AddDatabase(db2)

	conn := NewConn()
	conn.Pila = p

	request, err := http.NewRequest("DELETE",
		fmt.Sprintf("/databases/%s",
			db1.Name),
		nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()

	databaseHandle := conn.databaseHandler(db1.Name)
	databaseHandle.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusNoContent)
	}

	if len(conn.Pila.Databases) != 1 {
		t.Errorf("got %d database, expected %d", len(conn.Pila.Databases), 1)
	}
}

func TestDatabaseHandler_Gone(t *testing.T) {
	conn := NewConn()

	request, err := http.NewRequest("GET",
		fmt.Sprintf("/databases/%s",
			"nodb"),
		nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()

	databaseHandle := conn.databaseHandler(uuid.UUID("nodb").String())
	databaseHandle.ServeHTTP(response, request)

	if response.Code != http.StatusGone {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusGone)
	}
}

func TestStacksHandler_GET(t *testing.T) {
	s1 := pila.NewStack("stack1")
	s1.Push("foo")

	s2 := pila.NewStack("stack2")
	s2.Push(1)
	s2.Push(8)

	db := pila.NewDatabase("db")
	_ = db.AddStack(s1)
	_ = db.AddStack(s2)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	request, err := http.NewRequest("GET", "/databases/db/stacks", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	stacksHandle := conn.stacksHandler(db.ID.String())
	stacksHandle.ServeHTTP(response, request)

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
	}

	if response.Code != http.StatusOK {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusOK)
	}

	stacks, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if expected := `{"stacks":[{"id":"f0306fec639bd57fc2929c8b897b9b37","name":"stack1","peek":"foo","size":1},{"id":"dde8f895aea2ffa5546336146b9384e7","name":"stack2","peek":8,"size":2}]}`; string(stacks) != expected {
		t.Errorf("stacks are %s, expected %s", string(stacks), expected)
	}
}

func TestStacksHandler_GET_Name(t *testing.T) {
	s1 := pila.NewStack("stack1")
	s1.Push("bar")

	s2 := pila.NewStack("stack2")
	s2.Push(`{"a":"b"}`)

	db := pila.NewDatabase("db")
	_ = db.AddStack(s1)
	_ = db.AddStack(s2)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	request, err := http.NewRequest("GET", "/databases/db/stacks", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	stacksHandle := conn.stacksHandler("db")
	stacksHandle.ServeHTTP(response, request)

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
	}

	if response.Code != http.StatusOK {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusOK)
	}

	stacks, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if expected := `{"stacks":[{"id":"f0306fec639bd57fc2929c8b897b9b37","name":"stack1","peek":"bar","size":1},{"id":"dde8f895aea2ffa5546336146b9384e7","name":"stack2","peek":"{\"a\":\"b\"}","size":1}]}`; string(stacks) != expected {
		t.Errorf("stacks are %s, expected %s", string(stacks), expected)
	}
}

func TestStacksHandler_GET_Gone(t *testing.T) {
	db := pila.NewDatabase("db")

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	request, err := http.NewRequest("GET", "/databases/nodb/stacks", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	stacksHandle := conn.stacksHandler("nodb")
	stacksHandle.ServeHTTP(response, request)

	if response.Code != http.StatusGone {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusGone)
	}
}

func TestStacksHandler_GET_BadRequest(t *testing.T) {
	ch := make(chan int)

	stack := pila.NewStack("test-stack-channel")
	stack.Push(ch)

	db := pila.NewDatabase("db")
	_ = db.AddStack(stack)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	request, err := http.NewRequest("GET", "/databases/db/stacks", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	stacksHandle := conn.stacksHandler("db")
	stacksHandle.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusGone)
	}
}

func TestStacksHandler_PUT(t *testing.T) {
	db := pila.NewDatabase("db")

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	path := fmt.Sprintf("/databases/%s/stacks/?name=test-stack", db.ID.String())
	request, err := http.NewRequest("PUT", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	stacksHandle := conn.stacksHandler(db.ID.String())
	stacksHandle.ServeHTTP(response, request)

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
	}

	if response.Code != http.StatusCreated {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusCreated)
	}

	stack, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	expectedStack := `{"id":"bb4dabeeaa6e90108583ddbf49649427","name":"test-stack","peek":null,"size":0}`

	if string(stack) != expectedStack {
		t.Errorf("stack is %s, expected %s", string(stack), expectedStack)
	}
}

func TestStacksHandler_PUT_Name(t *testing.T) {
	db := pila.NewDatabase("db")

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	path := fmt.Sprintf("/databases/%s/stacks/?name=test-stack", db.Name)
	request, err := http.NewRequest("PUT", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	stacksHandle := conn.stacksHandler(db.ID.String())
	stacksHandle.ServeHTTP(response, request)

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
	}

	if response.Code != http.StatusCreated {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusCreated)
	}

	stack, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	expectedStack := `{"id":"bb4dabeeaa6e90108583ddbf49649427","name":"test-stack","peek":null,"size":0}`

	if string(stack) != expectedStack {
		t.Errorf("stack is %s, expected %s", string(stack), expectedStack)
	}
}

func TestCreateStackHandler(t *testing.T) {
	db := pila.NewDatabase("db")

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	path := fmt.Sprintf("/databases/%s/stacks/?name=test-stack", db.ID.String())
	request, err := http.NewRequest("PUT", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	conn.createStackHandler(response, request, db.ID.String())

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
	}

	if response.Code != http.StatusCreated {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusCreated)
	}

	stack, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	expectedStack := `{"id":"bb4dabeeaa6e90108583ddbf49649427","name":"test-stack","peek":null,"size":0}`
	if string(stack) != expectedStack {
		t.Errorf("stack is %s, expected %s", string(stack), expectedStack)
	}
}

func TestCreateStackHandler_NoName(t *testing.T) {
	db := pila.NewDatabase("db")

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	path := fmt.Sprintf("/databases/%s/stacks/?name=", db.ID.String())
	request, err := http.NewRequest("PUT", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	conn.createStackHandler(response, request, db.ID.String())

	if response.Code != http.StatusBadRequest {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusBadRequest)
	}
}

func TestCreateStackHandler_Gone(t *testing.T) {
	p := pila.NewPila()

	conn := NewConn()
	conn.Pila = p

	path := fmt.Sprintf("/databases/%s/stacks/?name=test-stack", "12345")
	request, err := http.NewRequest("PUT", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	conn.createStackHandler(response, request, "12345")

	if response.Code != http.StatusGone {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusGone)
	}
}

func TestCreateStackHandler_Conflict(t *testing.T) {
	s := pila.NewStack("test-stack")

	db := pila.NewDatabase("db")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	path := fmt.Sprintf("/databases/%s/stacks/?name=test-stack", db.ID.String())
	request, err := http.NewRequest("PUT", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	conn.createStackHandler(response, request, db.ID.String())

	if response.Code != http.StatusConflict {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusConflict)
	}
}

func TestStackHandler_POST(t *testing.T) {
	s := pila.NewStack("stack")

	db := pila.NewDatabase("db")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	element := pila.Element{Value: "test-element"}
	expectedElementJSON, _ := element.ToJSON()

	request, err := http.NewRequest("POST",
		fmt.Sprintf("/databases/%s/stacks/%s",
			db.ID.String(),
			s.ID.String()),
		bytes.NewBuffer(expectedElementJSON))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	params := &map[string]string{
		"database_id": db.ID.String(),
		"stack_id":    s.ID.String(),
	}

	stackHandle := conn.stackHandler(params)
	stackHandle.ServeHTTP(response, request)

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
	}

	if response.Code != http.StatusOK {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusOK)
	}

	elementJSON, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(elementJSON) != string(expectedElementJSON) {
		t.Errorf("pushed element is %v, expected %v", string(elementJSON), string(expectedElementJSON))
	}
}

func TestStackHandler_POST_Name(t *testing.T) {
	s := pila.NewStack("stack")

	db := pila.NewDatabase("db")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	element := pila.Element{Value: "test-element"}
	expectedElementJSON, _ := element.ToJSON()

	request, err := http.NewRequest("POST",
		fmt.Sprintf("/databases/%s/stacks/%s",
			db.Name,
			s.Name),
		bytes.NewBuffer(expectedElementJSON))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	params := &map[string]string{
		"database_id": db.ID.String(),
		"stack_id":    s.ID.String(),
	}

	stackHandle := conn.stackHandler(params)
	stackHandle.ServeHTTP(response, request)

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
	}

	if response.Code != http.StatusOK {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusOK)
	}

	elementJSON, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(elementJSON) != string(expectedElementJSON) {
		t.Errorf("pushed element is %v, expected %v", string(elementJSON), string(expectedElementJSON))
	}
}

func TestPushStackHandler(t *testing.T) {
	s := pila.NewStack("stack")

	db := pila.NewDatabase("db")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	element := pila.Element{Value: "test-element"}
	expectedElementJSON, _ := element.ToJSON()

	request, err := http.NewRequest("POST",
		fmt.Sprintf("/databases/%s/stacks/%s",
			db.ID.String(),
			s.ID.String()),
		bytes.NewBuffer(expectedElementJSON))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	vars := map[string]string{
		"database_id": db.ID.String(),
		"stack_id":    s.ID.String(),
	}

	conn.pushStackHandler(response, request, vars)

	if pushedElement := db.Stacks[s.ID].Peek(); pushedElement != element.Value {
		t.Errorf("Pushed element is %v, expected %v", pushedElement, element.Value)
	}

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
	}

	if response.Code != http.StatusOK {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusOK)
	}

	elementJSON, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(elementJSON) != string(expectedElementJSON) {
		t.Errorf("pushed element is %v, expected %v", string(elementJSON), string(expectedElementJSON))
	}
}

func TestPushStackHandler_Name(t *testing.T) {
	s := pila.NewStack("stack")

	db := pila.NewDatabase("db")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	element := pila.Element{Value: "test-element"}
	expectedElementJSON, _ := element.ToJSON()

	request, err := http.NewRequest("POST",
		fmt.Sprintf("/databases/%s/stacks/%s",
			db.ID.String(),
			s.ID.String()),
		bytes.NewBuffer(expectedElementJSON))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	vars := map[string]string{
		"database_id": db.Name,
		"stack_id":    s.Name,
	}

	conn.pushStackHandler(response, request, vars)

	if pushedElement := db.Stacks[s.ID].Peek(); pushedElement != element.Value {
		t.Errorf("Pushed element is %v, expected %v", pushedElement, element.Value)
	}

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
	}

	if response.Code != http.StatusOK {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusOK)
	}

	elementJSON, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(elementJSON) != string(expectedElementJSON) {
		t.Errorf("pushed element is %v, expected %v", string(elementJSON), string(expectedElementJSON))
	}
}

func TestPushStackHandler_Empty(t *testing.T) {
	s := pila.NewStack("stack")

	db := pila.NewDatabase("db")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	request, err := http.NewRequest("POST",
		fmt.Sprintf("/databases/%s/stacks/%s",
			db.ID.String(),
			s.ID.String()),
		nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	vars := map[string]string{
		"database_id": db.ID.String(),
		"stack_id":    s.ID.String(),
	}

	conn.pushStackHandler(response, request, vars)

	if pushedElement := db.Stacks[s.ID].Peek(); pushedElement != nil {
		t.Errorf("Pushed element is %v, expected nil", pushedElement)
	}

	if response.Code != http.StatusBadRequest {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusBadRequest)
	}
}

func TestPushStackHandler_DatabaseGone(t *testing.T) {
	p := pila.NewPila()

	conn := NewConn()
	conn.Pila = p

	element := pila.Element{Value: "test-element"}
	expectedElementJSON, _ := element.ToJSON()

	request, err := http.NewRequest("POST",
		fmt.Sprintf("/databases/%s/stacks/%s",
			"db",
			"stack"),
		bytes.NewBuffer(expectedElementJSON))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	vars := map[string]string{
		"database_id": "db",
		"stack_id":    "stack",
	}

	conn.pushStackHandler(response, request, vars)

	if response.Code != http.StatusGone {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusGone)
	}
}

func TestPushStackHandler_StackGone(t *testing.T) {
	p := pila.NewPila()

	db := pila.NewDatabase("db")
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	element := pila.Element{Value: "test-element"}
	expectedElementJSON, _ := element.ToJSON()

	request, err := http.NewRequest("POST",
		fmt.Sprintf("/databases/%s/stacks/%s",
			db.ID.String(),
			"stack"),
		bytes.NewBuffer(expectedElementJSON))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	vars := map[string]string{
		"database_id": db.ID.String(),
		"stack_id":    "stack",
	}

	conn.pushStackHandler(response, request, vars)

	if response.Code != http.StatusGone {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusGone)
	}
}

func TestPushStackHandler_BadDecoding(t *testing.T) {
	s := pila.NewStack("stack")

	db := pila.NewDatabase("db")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	element := []byte(`{`)

	request, err := http.NewRequest("POST",
		fmt.Sprintf("/databases/%s/stacks/%s",
			db.ID.String(),
			s.ID.String()),
		bytes.NewBuffer(element))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	vars := map[string]string{
		"database_id": db.ID.String(),
		"stack_id":    s.ID.String(),
	}

	conn.pushStackHandler(response, request, vars)

	if pushedElement := db.Stacks[s.ID].Peek(); pushedElement != nil {
		t.Errorf("Pushed element is %v, expected nil", pushedElement)
	}

	if response.Code != http.StatusBadRequest {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusBadRequest)
	}
}

func TestPopStackHandler(t *testing.T) {
	s := pila.NewStack("stack")
	s.Push("foo")

	db := pila.NewDatabase("db")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	request, err := http.NewRequest("GET",
		fmt.Sprintf("/databases/%s/stacks/%s",
			db.ID.String(),
			s.ID.String()),
		nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()

	params := map[string]string{
		"database_id": db.ID.String(),
		"stack_id":    s.ID.String(),
	}

	popStackHandle := conn.popStackHandler(params)
	popStackHandle.ServeHTTP(response, request)

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
	}

	if response.Code != http.StatusOK {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusOK)
	}

	elementJSON, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(elementJSON) != "\"foo\"" {
		t.Errorf("popped element is %v, expected %v", string(elementJSON), "\"foo\"")
	}
}

func TestPopStackHandler_EmptyStack(t *testing.T) {
	s := pila.NewStack("stack")

	db := pila.NewDatabase("db")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	request, err := http.NewRequest("GET",
		fmt.Sprintf("/databases/%s/stacks/%s",
			db.ID.String(),
			s.ID.String()),
		nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()

	params := map[string]string{
		"database_id": db.ID.String(),
		"stack_id":    s.ID.String(),
	}

	popStackHandle := conn.popStackHandler(params)
	popStackHandle.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusNoContent)
	}
}

func TestPopStackHandler_NoStackFound(t *testing.T) {
	s := pila.NewStack("stack")

	db := pila.NewDatabase("db")

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	request, err := http.NewRequest("GET",
		fmt.Sprintf("/databases/%s/stacks/%s",
			db.ID.String(),
			s.ID.String()),
		nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()

	params := map[string]string{
		"database_id": db.ID.String(),
		"stack_id":    s.ID.String(),
	}

	popStackHandle := conn.popStackHandler(params)
	popStackHandle.ServeHTTP(response, request)

	if response.Code != http.StatusGone {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusGone)
	}
}

func TestPopStackHandler_NoDatabaseFound(t *testing.T) {
	s := pila.NewStack("stack")
	db := pila.NewDatabase("db")

	p := pila.NewPila()

	conn := NewConn()
	conn.Pila = p

	request, err := http.NewRequest("GET",
		fmt.Sprintf("/databases/%s/stacks/%s",
			db.ID.String(),
			s.ID.String()),
		nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()

	params := map[string]string{
		"database_id": db.ID.String(),
		"stack_id":    s.ID.String(),
	}

	popStackHandle := conn.popStackHandler(params)
	popStackHandle.ServeHTTP(response, request)

	if response.Code != http.StatusGone {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusGone)
	}
}

func TestNotFoundHandler_WrongEndpoint(t *testing.T) {
	conn := NewConn()
	request, err := http.NewRequest("GET", "/_statuss", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	conn.notFoundHandler(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusNotFound)
	}
}

func TestNotFoundHandler_WrongType(t *testing.T) {
	conn := NewConn()
	request, err := http.NewRequest("POST", "/_status", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	conn.notFoundHandler(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusNotFound)
	}
}

func TestGoneHandler(t *testing.T) {
	conn := NewConn()
	request, err := http.NewRequest("GET", "/databases/nodb", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	conn.goneHandler(response, request, "database nodb is Gone")

	if response.Code != http.StatusGone {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusNotFound)
	}
}
