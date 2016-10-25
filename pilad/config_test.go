package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fern4lvarez/piladb/config"
	"github.com/fern4lvarez/piladb/config/vars"
	"github.com/fern4lvarez/piladb/pila"
)

func TestConfigHandler_GET(t *testing.T) {
	db := pila.NewDatabase("db")

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p
	conn.Config = config.NewConfig()
	conn.Config.Set("SIZE", 2)
	conn.Config.Set("PORT", "1205")
	conn.Config.Set("PORT", "8080")

	inputOutput := []struct {
		input  string
		output string
	}{
		{input: "/_config", output: `{"stacks":{"PORT":"8080","SIZE":2}}`},
	}

	for _, io := range inputOutput {
		request, err := http.NewRequest("GET", io.input, nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()

		conn.configHandler(response, request)

		if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
			t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
		}

		if response.Code != http.StatusOK {
			t.Errorf("response code is %v, expected %v", response.Code, http.StatusOK)
		}

		config, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		if string(config) != io.output {
			t.Errorf("config is %s, expected %s", string(config), io.output)
		}
	}
}

func TestConfigHandler_GET_BadRequest(t *testing.T) {
	ch := make(chan int)

	db := pila.NewDatabase("db")

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p
	conn.Config = config.NewConfig()
	conn.Config.Set("SIZE", ch)

	request, err := http.NewRequest("GET", "/_config", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	conn.configHandler(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusBadRequest)
	}
}

func TestCheckMaxSizeOfStack(t *testing.T) {
	s := pila.NewStack("stack")
	s.Push("foo")

	db := pila.NewDatabase("mydb")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	f := func(w http.ResponseWriter, r *http.Request, stack *pila.Stack) {
		w.WriteHeader(http.StatusOK)
	}

	inputOutput := []struct {
		input  int
		output int
	}{
		{input: 1, output: http.StatusNotAcceptable},
		{input: 6, output: http.StatusOK},
	}

	for _, io := range inputOutput {
		conn.Config.Set(vars.MaxSizeOfStack, io.input)
		request, err := http.NewRequest("GET", "", nil)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		conn.checkMaxSizeOfStack(f)(response, request, s)

		if response.Code != io.output {
			t.Errorf("response code is %v, expected %v", response.Code, io.output)
		}
	}
}
