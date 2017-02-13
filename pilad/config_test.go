package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/fern4lvarez/piladb/config"
	"github.com/fern4lvarez/piladb/config/vars"
	"github.com/fern4lvarez/piladb/pila"
)

func TestBuildConfig(t *testing.T) {
	conn := NewConn()

	if err := os.Unsetenv(vars.Env(vars.MaxStackSize)); err != nil {
		t.Fatal(err)
	}

	maxStackSizeFlag = 32
	conn.buildConfig()

	if s := conn.Config.Get(vars.MaxStackSize); s != 32 {
		t.Errorf("MaxStackSize is %v, expected %d", s, 32)
	}

	if err := os.Setenv(vars.Env(vars.MaxStackSize), "42"); err != nil {
		t.Fatal(err)
	}
	conn.buildConfig()

	if s := conn.Config.Get(vars.MaxStackSize); s != 42 {
		t.Errorf("MaxStackSize is %v, expected %d", s, 42)
	}

	if err := os.Setenv(vars.Env(vars.MaxStackSize), "foo"); err != nil {
		t.Fatal(err)
	}
	conn.buildConfig()

	if s := conn.Config.Get(vars.MaxStackSize); s != -1 {
		t.Errorf("MaxStackSize is %v, expected %v", s, -1)
	}

	pushWhenFullFlag = true
	if err := os.Setenv(vars.Env(vars.PushWhenFull), "false"); err != nil {
		t.Fatal(err)
	}
	conn.buildConfig()

	if s := conn.Config.PushWhenFull(); s != false {
		t.Errorf("PushWhenFull is %v, expected %v", s, false)
	}

	if err := os.Setenv(vars.Env(vars.PushWhenFull), "true"); err != nil {
		t.Fatal(err)
	}
	conn.buildConfig()

	if s := conn.Config.PushWhenFull(); s != true {
		t.Errorf("PushWhenFull is %v, expected %v", s, true)
	}
}

func TestConfigHandler_GET(t *testing.T) {
	conn := NewConn()
	conn.Config = config.NewConfig()
	conn.Config.Set("SIZE", 2)
	conn.Config.Set("PORT", "1205")
	conn.Config.Set("PORT", "8080")

	inputOutput := []struct {
		input, output string
	}{
		{"/_config", `{"stacks":{"PORT":"8080","SIZE":2}}`},
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

	conn := NewConn()
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

func TestConfigKeyHandler(t *testing.T) {
	conn := NewConn()
	conn.Config = config.NewConfig()
	conn.Config.Set(vars.MaxStackSize, 2)

	element := pila.Element{Value: 2}
	expectedElementJSON, _ := element.ToJSON()

	newElement := pila.Element{Value: 10}
	expectedNewElementJSON, _ := newElement.ToJSON()

	inputOutput := []struct {
		input struct {
			method, key string
			payload     io.Reader
		}
		output struct {
			value    interface{}
			response []byte
		}
	}{
		{struct {
			method, key string
			payload     io.Reader
		}{"GET", vars.MaxStackSize, nil},
			struct {
				value    interface{}
				response []byte
			}{2, expectedElementJSON},
		},
		{struct {
			method, key string
			payload     io.Reader
		}{"POST", vars.MaxStackSize, bytes.NewBuffer(expectedNewElementJSON)},
			struct {
				value    interface{}
				response []byte
			}{10, expectedNewElementJSON},
		},
	}

	for _, io := range inputOutput {
		request, err := http.NewRequest(io.input.method,
			fmt.Sprintf("/_config/%s", io.input.key),
			io.input.payload)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		configKeyHandle := conn.configKeyHandler(io.input.key)
		configKeyHandle.ServeHTTP(response, request)

		if value := conn.Config.MaxStackSize(); value != io.output.value {
			t.Errorf("value is %d, expected %d", value, io.output.value)
		}

		if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
			t.Errorf("Content-Type is %v, expected %v", contentType, "application/json")
		}

		if response.Code != http.StatusOK {
			t.Errorf("response code is %v, expected %v", response.Code, http.StatusOK)
		}

		responseJSON, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		if string(responseJSON) != string(io.output.response) {
			t.Errorf("response is %s, expected %s", string(responseJSON), string(io.output.response))
		}
	}
}

func TestConfigKeyHandler_Gone(t *testing.T) {
	key := "no-exist"

	conn := NewConn()
	conn.Config = config.NewConfig()

	request, err := http.NewRequest("GET", fmt.Sprintf("/_config/%s", key), nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	configKeyHandle := conn.configKeyHandler(key)
	configKeyHandle.ServeHTTP(response, request)

	if response.Code != http.StatusGone {
		t.Errorf("response code is %v, expected %v", response.Code, http.StatusGone)
	}
}

func TestConfigKeyHandler_BadRequest(t *testing.T) {
	ch := make(chan int)
	key := "SIZE"

	conn := NewConn()
	conn.Config = config.NewConfig()
	conn.Config.Set(key, ch)

	input := []struct {
		method  string
		payload io.Reader
	}{
		{"GET", nil},
		{"POST", nil},
		{"POST", bytes.NewBuffer([]byte("{"))},
	}

	for _, in := range input {
		request, err := http.NewRequest(in.method,
			fmt.Sprintf("/_config/%s", key), in.payload)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()

		configKeyHandle := conn.configKeyHandler(key)
		configKeyHandle.ServeHTTP(response, request)

		if response.Code != http.StatusBadRequest {
			t.Errorf("response code is %v, expected %v", response.Code, http.StatusBadRequest)
		}
	}
}

func TestCheckMaxStackSize(t *testing.T) {
	s := pila.NewStack("stack", time.Now())
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
		input, output int
	}{
		{1, http.StatusNotAcceptable},
		{6, http.StatusOK},
	}

	for _, io := range inputOutput {
		conn.Config.Set(vars.MaxStackSize, io.input)
		request, err := http.NewRequest("POST",
			fmt.Sprintf("/databases/%s/stacks/%s",
				db.ID.String(),
				s.ID.String()),
			nil)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		conn.checkMaxStackSize(f)(response, request, s)

		if response.Code != io.output {
			t.Errorf("response code is %v, expected %v", response.Code, io.output)
		}
	}
}

func TestCheckMaxStackSize_PushWhenFull(t *testing.T) {
	s := pila.NewStack("stack", time.Now())
	s.Push("foo")

	db := pila.NewDatabase("mydb")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	conn.Config.Set(vars.PushWhenFull, true)

	element := pila.Element{Value: "test-element"}
	elementJSON, _ := element.ToJSON()

	f := func(w http.ResponseWriter, r *http.Request, stack *pila.Stack) {
		w.WriteHeader(http.StatusOK)
	}

	inputOutput := []struct {
		inputSize                int
		inputPush                bool
		outputStatus, outputSize int
		outputPeek, outputSweep  interface{}
	}{
		{1, false, http.StatusNotAcceptable, 1, "foo", nil},
		{1, false, http.StatusNotAcceptable, 1, "foo", nil},
		{1, true, http.StatusOK, 1, "foo", true},
	}

	for _, io := range inputOutput {
		conn.Config.Set(vars.MaxStackSize, io.inputSize)
		conn.Config.Set(vars.PushWhenFull, io.inputPush)
		request, err := http.NewRequest("POST",
			fmt.Sprintf("/databases/%s/stacks/%s",
				db.ID.String(),
				s.ID.String()),
			bytes.NewBuffer(elementJSON))
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		conn.checkMaxStackSize(f)(response, request, s)

		if response.Code != io.outputStatus {
			t.Errorf("response code is %v, expected %v", response.Code, io.outputStatus)
		}

		if s.Size() != io.outputSize {
			t.Errorf("Stack size is %v, expected %v", s.Size(), io.outputSize)
		}
		if s.Peek() != io.outputPeek {
			t.Errorf("Stack peek is %v, expected %v", s.Peek(), io.outputPeek)
		}

		if sweepBeforePush := conn.Config.Get("SWEEP_BEFORE_PUSH"); sweepBeforePush != io.outputSweep {
			t.Errorf("SWEEP_BEFORE_PUSH is %v, expected %v", sweepBeforePush, io.outputSweep)
		}
	}
}

func TestCheckMaxStackSize_PushWhenFullWithStackEmpty(t *testing.T) {
	s := pila.NewStack("stack", time.Now())

	db := pila.NewDatabase("mydb")
	_ = db.AddStack(s)

	p := pila.NewPila()
	_ = p.AddDatabase(db)

	conn := NewConn()
	conn.Pila = p

	conn.Config.Set(vars.PushWhenFull, true)

	element := pila.Element{Value: "test-element"}
	elementJSON, _ := element.ToJSON()

	f := func(w http.ResponseWriter, r *http.Request, stack *pila.Stack) {
		w.WriteHeader(http.StatusOK)
	}

	inputOutput := []struct {
		inputSize                int
		inputPush                bool
		outputStatus, outputSize int
		outputPeek               interface{}
	}{
		{0, true, http.StatusOK, 0, nil},
	}

	for _, io := range inputOutput {
		conn.Config.Set(vars.MaxStackSize, io.inputSize)
		conn.Config.Set(vars.PushWhenFull, io.inputPush)
		request, err := http.NewRequest("POST",
			fmt.Sprintf("/databases/%s/stacks/%s",
				db.ID.String(),
				s.ID.String()),
			bytes.NewBuffer(elementJSON))
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		conn.checkMaxStackSize(f)(response, request, s)

		if response.Code != io.outputStatus {
			t.Errorf("response code is %v, expected %v", response.Code, io.outputStatus)
		}

		if s.Size() != io.outputSize {
			t.Errorf("Stack size is %v, expected %v", s.Size(), io.outputSize)
		}
		if s.Peek() != io.outputPeek {
			t.Errorf("Stack peek is %v, expected %v", s.Peek(), io.outputPeek)
		}
	}
}
