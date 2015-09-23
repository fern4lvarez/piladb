package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

	if response.Code != 200 {
		t.Errorf("response code is %v, expected %v", response.Code, 200)
	}

	statusJSON, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(string(statusJSON), `{"status":"OK","version"`) {
		t.Errorf("status is %s", string(statusJSON))
	}
}

func TestDatabasesHandler(t *testing.T) {
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

	if response.Code != 200 {
		t.Errorf("response code is %v, expected %v", response.Code, 200)
	}

	databases, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(databases) != `{"number_of_databases":0,"databases":[]}` {
		t.Errorf("databases are %s, expected %s", string(databases), `{"number_of_databases":0,"databases":[]}`)
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

	if response.Code != 404 {
		t.Errorf("response code is %v, expected %v", response.Code, 404)
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

	if response.Code != 404 {
		t.Errorf("response code is %v, expected %v", response.Code, 404)
	}
}
