package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/fern4lvarez/piladb/pila"
	"github.com/fern4lvarez/piladb/pkg/date"
	"github.com/fern4lvarez/piladb/pkg/uuid"
	"github.com/fern4lvarez/piladb/pkg/version"
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

func TestRootHandler(t *testing.T) {
	conn := NewConn()
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	expectedRedirAddress := fmt.Sprintf("https://raw.githubusercontent.com/fern4lvarez/piladb/%s/pilad/README.md", version.CommitHash())

	conn.rootHandler(response, request)

	if response.Code != http.StatusMovedPermanently {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusMovedPermanently)
	}

	locations, ok := response.Header()["Location"]
	if !ok {
		t.Fatal("no Location Header found")
	}

	if l := len(locations); l != 1 {
		t.Fatalf("number of redirections is %d, expected %d", l, 1)
	}

	if l := locations[0]; l != expectedRedirAddress {
		t.Errorf("redirection Address is %s, expected %s", l, expectedRedirAddress)
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
	s := pila.NewStack("stack", time.Now().UTC())
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
	s := pila.NewStack("stack", time.Now().UTC())
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
	now1 := time.Now().UTC()
	after1 := time.Now().UTC()
	now2 := time.Now().UTC()
	after2 := time.Now().UTC()

	s1 := pila.NewStack("stack1", now1)
	s1.Push("foo")
	s1.Update(after1)

	s2 := pila.NewStack("stack2", now2)
	s2.Push(1)
	s2.Push(8)
	s2.Update(after2)

	db := pila.NewDatabase("db")
	_ = db.AddStack(s1)
	_ = db.AddStack(s2)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	inputOutput := []struct {
		input, output string
	}{
		{"/databases/db/stacks", fmt.Sprintf(`{"stacks":[{"id":"f0306fec639bd57fc2929c8b897b9b37","name":"stack1","peek":"foo","size":1,"created_at":"%v","updated_at":"%v","read_at":"%v"},{"id":"dde8f895aea2ffa5546336146b9384e7","name":"stack2","peek":8,"size":2,"created_at":"%v","updated_at":"%v","read_at":"%v"}]}`,
			date.Format(now1.Local()), date.Format(after1.Local()), date.Format(after1.Local()),
			date.Format(now2.Local()), date.Format(after2.Local()), date.Format(after2.Local()))},
		{"/databases/db/stacks?kv", `{"stacks":{"stack1":"foo","stack2":8}}`},
	}

	for _, io := range inputOutput {
		request, err := http.NewRequest("GET", io.input, nil)
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

		if string(stacks) != io.output {
			t.Errorf("stacks are %s, expected %s", string(stacks), io.output)
		}
	}
}

func TestStacksHandler_GET_Name(t *testing.T) {
	now1 := time.Now().UTC()
	after1 := time.Now().UTC()
	now2 := time.Now().UTC()
	after2 := time.Now().UTC()

	s1 := pila.NewStack("stack1", now1)
	s1.Push("bar")
	s1.Update(after1)

	s2 := pila.NewStack("stack2", now2)
	s2.Push(`{"a":"b"}`)
	s2.Update(after2)

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

	if expected := fmt.Sprintf(`{"stacks":[{"id":"f0306fec639bd57fc2929c8b897b9b37","name":"stack1","peek":"bar","size":1,"created_at":"%v","updated_at":"%v","read_at":"%v"},{"id":"dde8f895aea2ffa5546336146b9384e7","name":"stack2","peek":"{\"a\":\"b\"}","size":1,"created_at":"%v","updated_at":"%v","read_at":"%v"}]}`,
		date.Format(now1.Local()), date.Format(after1.Local()), date.Format(after1.Local()),
		date.Format(now2.Local()), date.Format(after2.Local()), date.Format(after2.Local())); string(stacks) != expected {
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

	stack := pila.NewStack("test-stack-channel", time.Now().UTC())
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
	conn.opDate = time.Now().UTC()

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

	expectedStack := fmt.Sprintf(`{"id":"bb4dabeeaa6e90108583ddbf49649427","name":"test-stack","peek":null,"size":0,"created_at":"%v","updated_at":"%v","read_at":"%v"}`,
		date.Format(conn.opDate.Local()), date.Format(conn.opDate.Local()), date.Format(conn.opDate.Local()))

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
	conn.opDate = time.Now().UTC()

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

	expectedStack := fmt.Sprintf(`{"id":"bb4dabeeaa6e90108583ddbf49649427","name":"test-stack","peek":null,"size":0,"created_at":"%v","updated_at":"%v","read_at":"%v"}`,
		date.Format(conn.opDate.Local()), date.Format(conn.opDate.Local()), date.Format(conn.opDate.Local()))

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
	conn.opDate = time.Now().UTC()

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

	expectedStack := fmt.Sprintf(`{"id":"bb4dabeeaa6e90108583ddbf49649427","name":"test-stack","peek":null,"size":0,"created_at":"%v","updated_at":"%v","read_at":"%v"}`,
		date.Format(conn.opDate.Local()), date.Format(conn.opDate.Local()), date.Format(conn.opDate.Local()))
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
	s := pila.NewStack("test-stack", time.Now().UTC())

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

func TestStackHandler_GET(t *testing.T) {
	element := pila.Element{Value: "test-element"}
	expectedElementJSON, _ := element.ToJSON()

	createDate := time.Now().UTC()
	s := pila.NewStack("stack", createDate)
	s.Update(createDate)

	db := pila.NewDatabase("db")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p
	conn.opDate = time.Now().UTC()

	s.Push(element.Value)

	expectedStackStatusJSON, err := s.Status().ToJSON()
	if err != nil {
		t.Fatal(err)
	}

	expectedSizeJSON := s.SizeToJSON()

	inputOutput := []struct {
		input struct {
			database, stack, op string
		}
		output struct {
			response []byte
			code     int
		}
	}{
		{struct {
			database, stack, op string
		}{db.ID.String(), s.ID.String(), ""},
			struct {
				response []byte
				code     int
			}{expectedStackStatusJSON, http.StatusOK},
		},
		{struct {
			database, stack, op string
		}{db.ID.String(), s.ID.String(), "peek"},
			struct {
				response []byte
				code     int
			}{expectedElementJSON, http.StatusOK},
		},
		{struct {
			database, stack, op string
		}{db.ID.String(), s.ID.String(), "size"},
			struct {
				response []byte
				code     int
			}{expectedSizeJSON, http.StatusOK},
		},
	}

	for _, io := range inputOutput {
		request, err := http.NewRequest("GET",
			fmt.Sprintf("/databases/%s/stacks/%s?%s",
				io.input.database,
				io.input.stack,
				io.input.op),
			nil)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		params := map[string]string{
			"database_id": io.input.database,
			"stack_id":    io.input.stack,
		}

		stackHandle := conn.stackHandler(&params)
		stackHandle.ServeHTTP(response, request)

		if peek := db.Stacks[s.ID].Peek(); peek != element.Value {
			t.Errorf("peek is %v, expected %v", peek, element.Value)
		}

		if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
			t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
		}

		if response.Code != io.output.code {
			t.Errorf("on %s response code is %v, expected %v", io.input.op, response.Code, io.output.code)
		}

		responseJSON, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		if io.input.op == "" {
			expectedStackStatusJSON, err := s.Status().ToJSON()
			if err != nil {
				t.Fatal(err)
			} else if string(responseJSON) != string(expectedStackStatusJSON) {
				t.Errorf("on %s response is %s, expected %s", io.input.op, string(responseJSON), string(expectedStackStatusJSON))
			}

		} else if string(responseJSON) != string(io.output.response) {
			t.Errorf("on %s response is %s, expected %s", io.input.op, string(responseJSON), string(io.output.response))
		}
	}
}

func TestStackHandler_POST(t *testing.T) {
	s := pila.NewStack("stack", time.Now().UTC())

	db := pila.NewDatabase("db")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p
	conn.opDate = time.Now().UTC()

	element := pila.Element{Value: "test-element"}
	expectedElementJSON, _ := element.ToJSON()

	paramss := []map[string]string{
		{
			"database_id": db.ID.String(),
			"stack_id":    s.ID.String(),
		},
		{
			"database_id": db.Name,
			"stack_id":    s.Name,
		},
	}

	for _, params := range paramss {
		request, err := http.NewRequest("POST",
			fmt.Sprintf("/databases/%s/stacks/%s",
				params["database_id"],
				params["stack_id"]),
			bytes.NewBuffer(expectedElementJSON))
		if err != nil {
			t.Fatal(err)
		}
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()

		stackHandle := conn.stackHandler(&params)
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
}

func TestStackHandler_DELETE(t *testing.T) {
	element := pila.Element{Value: "test-element"}
	expectedElementJSON, _ := element.ToJSON()

	s := pila.NewStack("stack", time.Now().UTC())

	db := pila.NewDatabase("db")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	s.Push(element.Value)

	inputOutput := []struct {
		input struct {
			database, stack, op string
		}
		output struct {
			response []byte
			code     int
		}
	}{
		{struct {
			database, stack, op string
		}{db.ID.String(), s.ID.String(), ""},
			struct {
				response []byte
				code     int
			}{expectedElementJSON, http.StatusOK},
		},
		{struct {
			database, stack, op string
		}{db.Name, s.Name, ""},
			struct {
				response []byte
				code     int
			}{expectedElementJSON, http.StatusOK},
		},
		{struct {
			database, stack, op string
		}{db.Name, s.Name, "flush"},
			struct {
				response []byte
				code     int
			}{nil, http.StatusOK},
		},
		{struct {
			database, stack, op string
		}{db.Name, s.Name, "full"},
			struct {
				response []byte
				code     int
			}{nil, http.StatusNoContent},
		},
	}

	for _, io := range inputOutput {
		request, err := http.NewRequest("DELETE",
			fmt.Sprintf("/databases/%s/stacks/%s?%s",
				io.input.database,
				io.input.stack,
				io.input.op),
			nil)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		params := map[string]string{
			"database_id": io.input.database,
			"stack_id":    io.input.stack,
		}

		stackHandle := conn.stackHandler(&params)
		stackHandle.ServeHTTP(response, request)

		if response.Code != io.output.code {
			t.Errorf("on op %s response code is %v, expected %v", io.input.op, response.Code, io.output.code)
		}

		responseJSON, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		if io.input.op == "flush" {
			s.Update(conn.opDate)
			stackStatus := s.Status()

			expectedStackStatusJSON, err := stackStatus.ToJSON()
			if err != nil {
				t.Fatal(err)
			}
			if string(responseJSON) != string(expectedStackStatusJSON) {
				t.Errorf("on op %s response is %s, expected %s", io.input.op, string(responseJSON), string(expectedStackStatusJSON))
			}
		} else if string(responseJSON) != string(io.output.response) {
			t.Errorf("on op %s response is %s, expected %s", io.input.op, string(responseJSON), string(io.output.response))
		}

		if io.input.op == "full" {
			if s := db.Stacks[uuid.UUID(io.input.stack)]; s != nil {
				t.Errorf("db contains %v, expected not to", io.input.stack)
			}
		} else {
			if peek, ok := db.Stacks[s.ID].Pop(); ok {
				t.Errorf("stack contains %v, expected to be empty", peek)
			}

			if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
			}

			// restore element for next table test iteration
			s.Push(element.Value)
		}
	}
}

func TestStackHandler_DatabaseGone(t *testing.T) {
	now := time.Now().UTC()
	s := pila.NewStack("stack", now)

	db := pila.NewDatabase("db")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p
	conn.opDate = now

	request, err := http.NewRequest("GET",
		fmt.Sprintf("/databases/%s/stacks/%s",
			"non-existing-db",
			s.ID.String()),
		nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	params := &map[string]string{
		"database_id": "non-existing-db",
		"stack_id":    s.ID.String(),
	}

	stackHandle := conn.stackHandler(params)
	stackHandle.ServeHTTP(response, request)

	if response.Code != http.StatusGone {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusOK)
	}
}

func TestStackHandler_StackGone(t *testing.T) {
	now := time.Now().UTC()
	s := pila.NewStack("stack", now)

	db := pila.NewDatabase("db")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p
	conn.opDate = now

	request, err := http.NewRequest("GET",
		fmt.Sprintf("/databases/%s/stacks/%s",
			db.ID.String(),
			"non-existing-stack"),
		nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	params := &map[string]string{
		"database_id": db.ID.String(),
		"stack_id":    "non-existing-stack",
	}

	stackHandle := conn.stackHandler(params)
	stackHandle.ServeHTTP(response, request)

	if response.Code != http.StatusGone {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusOK)
	}
}

func TestStatusStackHandler(t *testing.T) {
	s := pila.NewStack("stack", time.Now().UTC())

	db := pila.NewDatabase("db")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	s.Push("one")

	expectedStackStatusJSON, err := s.Status().ToJSON()
	if err != nil {
		t.Fatal(err)
	}

	varss := []map[string]string{
		{
			"database_id": db.ID.String(),
			"stack_id":    s.ID.String(),
		},
		{
			"database_id": db.Name,
			"stack_id":    s.Name,
		},
	}

	for _, vars := range varss {
		request, err := http.NewRequest("GET",
			fmt.Sprintf("/databases/%s/stacks/%s",
				vars["database_id"],
				vars["stack_id"]),
			nil)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		conn.statusStackHandler(response, request, s)
		if peek := db.Stacks[s.ID].Peek(); peek != "one" {
			t.Errorf("peek is %v, expected %v", peek, "one")
		}

		if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
			t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
		}

		if response.Code != http.StatusOK {
			t.Errorf("response code is %v, expected %v", response.Code, http.StatusOK)
		}

		stackStatusJSON, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		if string(stackStatusJSON) != string(expectedStackStatusJSON) {
			t.Errorf("stack status is %s, expected %s", string(stackStatusJSON), string(expectedStackStatusJSON))
		}
	}
}

func TestPeekStackHandler(t *testing.T) {
	s := pila.NewStack("stack", time.Now().UTC())

	db := pila.NewDatabase("db")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	element := pila.Element{Value: "test-element"}
	expectedElementJSON, _ := element.ToJSON()

	s.Push(element.Value)

	request, err := http.NewRequest("GET",
		fmt.Sprintf("/databases/%s/stacks/%s?peek",
			db.ID.String(),
			s.ID.String()),
		nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	conn.peekStackHandler(response, request, s)

	if peekElement := db.Stacks[s.ID].Peek(); peekElement != element.Value {
		t.Errorf("peek element is %v, expected %v", peekElement, element.Value)
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
		t.Errorf("peek element is %v, expected %v", string(elementJSON), string(expectedElementJSON))
	}
}

func TestSizeStackHandler(t *testing.T) {
	s := pila.NewStack("stack", time.Now().UTC())

	db := pila.NewDatabase("db")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	s.Push("element")

	expectedSizeJSON := s.SizeToJSON()

	request, err := http.NewRequest("GET",
		fmt.Sprintf("/databases/%s/stacks/%s?size",
			db.ID.String(),
			s.ID.String()),
		nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	conn.sizeStackHandler(response, request, s)

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
	}

	if response.Code != http.StatusOK {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusOK)
	}

	sizeJSON, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(sizeJSON) != string(expectedSizeJSON) {
		t.Errorf("size is %v, expected %v", string(sizeJSON), string(expectedSizeJSON))
	}
}

func TestPushStackHandler(t *testing.T) {
	s := pila.NewStack("stack", time.Now().UTC())

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

	conn.pushStackHandler(response, request, s)

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
	s := pila.NewStack("stack", time.Now().UTC())

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

	conn.pushStackHandler(response, request, s)

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

func TestPushStackHandler_SweepBefore(t *testing.T) {
	s := pila.NewStack("stack", time.Now().UTC())
	s.Push("foo")

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

	conn.Config.Set("SWEEP_BEFORE_PUSH", true)
	conn.pushStackHandler(response, request, s)

	if pushedElement := db.Stacks[s.ID].Peek(); pushedElement != element.Value {
		t.Errorf("Pushed element is %v, expected %v", pushedElement, element.Value)
	}

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
	}

	if response.Code != http.StatusOK {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusOK)
	}

	if size := db.Stacks[s.ID].Size(); size != 1 {
		t.Errorf("Stack size is %d, expected %d", size, 1)
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
	s := pila.NewStack("stack", time.Now().UTC())

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

	conn.pushStackHandler(response, request, s)

	if pushedElement := db.Stacks[s.ID].Peek(); pushedElement != nil {
		t.Errorf("Pushed element is %v, expected nil", pushedElement)
	}

	if response.Code != http.StatusBadRequest {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusBadRequest)
	}
}

func TestPushStackHandler_BadDecoding(t *testing.T) {
	s := pila.NewStack("stack", time.Now().UTC())

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

	conn.pushStackHandler(response, request, s)

	if pushedElement := db.Stacks[s.ID].Peek(); pushedElement != nil {
		t.Errorf("Pushed element is %v, expected nil", pushedElement)
	}

	if response.Code != http.StatusBadRequest {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusBadRequest)
	}
}

func TestPopStackHandler(t *testing.T) {
	element := pila.Element{Value: "test-element"}
	expectedElementJSON, _ := element.ToJSON()

	s := pila.NewStack("stack", time.Now().UTC())
	s.Push(element.Value)

	db := pila.NewDatabase("db")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	varss := []map[string]string{
		{
			"database_id": db.ID.String(),
			"stack_id":    s.ID.String(),
		},
		{
			"database_id": db.Name,
			"stack_id":    s.Name,
		},
	}

	for _, vars := range varss {
		request, err := http.NewRequest("DELETE",
			fmt.Sprintf("/databases/%s/stacks/%s",
				vars["database_id"],
				vars["stack_id"]),
			nil)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		conn.popStackHandler(response, request, s)

		if peek, ok := db.Stacks[s.ID].Pop(); ok {
			t.Errorf("stack contains %v, expected to be empty", peek)
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
			t.Errorf("popped element is %s, expected %s", string(elementJSON), string(expectedElementJSON))
		}

		// restore element for next table test iteration
		s.Push(element.Value)
	}
}

func TestPopStackHandler_EmptyStack(t *testing.T) {
	s := pila.NewStack("stack", time.Now().UTC())

	db := pila.NewDatabase("db")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	request, err := http.NewRequest("DELETE",
		fmt.Sprintf("/databases/%s/stacks/%s",
			db.ID.String(),
			s.ID.String()),
		nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()

	conn.popStackHandler(response, request, s)

	if response.Code != http.StatusNoContent {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusNoContent)
	}
}

func TestFlushStackHandler(t *testing.T) {
	s := pila.NewStack("stack", time.Now().UTC())

	db := pila.NewDatabase("db")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	expectedStackStatusJSON, err := s.Status().ToJSON()
	if err != nil {
		t.Fatal(err)
	}

	s.Push("one")
	s.Push("two")
	s.Push("three")

	varss := []map[string]string{
		{
			"database_id": db.ID.String(),
			"stack_id":    s.ID.String(),
		},
		{
			"database_id": db.Name,
			"stack_id":    s.Name,
		},
	}

	for _, vars := range varss {
		request, err := http.NewRequest("DELETE",
			fmt.Sprintf("/databases/%s/stacks/%s?flush",
				vars["database_id"],
				vars["stack_id"]),
			nil)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		conn.flushStackHandler(response, request, s)

		if peek, ok := db.Stacks[s.ID].Pop(); ok {
			t.Errorf("stack contains %v, expected to be empty", peek)
		}

		if size := db.Stacks[s.ID].Size(); size != 0 {
			t.Errorf("stack has size %d, expected %d", size, 0)
		}

		if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
			t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
		}

		if response.Code != http.StatusOK {
			t.Errorf("response code is %v, expected %v", response.Code, http.StatusOK)
		}

		stackStatusJSON, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		if string(stackStatusJSON) != string(expectedStackStatusJSON) {
			t.Errorf("stack status is %s, expected %s", string(stackStatusJSON), string(expectedStackStatusJSON))
		}

		// restore elements for next table test iteration
		s.Push("one")
		s.Push("two")
		s.Push("three")
	}
}

func TestDeleteStackHandler(t *testing.T) {
	s := pila.NewStack("stack", time.Now().UTC())

	db := pila.NewDatabase("db")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	s.Push("one")

	varss := []map[string]string{
		{
			"database_id": db.ID.String(),
			"stack_id":    s.ID.String(),
		},
		{
			"database_id": db.Name,
			"stack_id":    s.Name,
		},
	}

	for _, vars := range varss {
		request, err := http.NewRequest("DELETE",
			fmt.Sprintf("/databases/%s/stacks/%s?full",
				vars["database_id"],
				vars["stack_id"]),
			nil)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		conn.deleteStackHandler(response, request, db, s)

		if expectedStack := db.Stacks[uuid.UUID(vars["stack_id"])]; expectedStack != nil {
			t.Errorf("db contains %v, expected not to", vars["stack_id"])
		}

		if response.Code != http.StatusNoContent {
			t.Errorf("response code is %v, expected %v", response.Code, http.StatusNoContent)
		}

		// restore elements for next table test iteration
		s = pila.NewStack("stack", time.Now().UTC())
		_ = db.AddStack(s)
		s.Push("one")
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
