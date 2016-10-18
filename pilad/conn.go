package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fern4lvarez/piladb/pila"
	"github.com/fern4lvarez/piladb/pkg/uuid"
	"github.com/fern4lvarez/piladb/pkg/version"

	"github.com/gorilla/mux"
)

// Conn represents the current piladb connection, containing
// the Pila instance and its status.
type Conn struct {
	Pila   *pila.Pila
	Status *Status
}

// NewConn creates and returns a new piladb connection.
func NewConn() *Conn {
	conn := &Conn{}
	conn.Pila = pila.NewPila()
	conn.Status = NewStatus(version.CommitHash(), time.Now(), MemStats())
	return conn
}

// Connection Handlers

// statusHandler writes the piladb status into the response.
func (c *Conn) statusHandler(w http.ResponseWriter, r *http.Request) {
	c.Status.Update(time.Now(), MemStats())

	w.Header().Set("Content-Type", "application/json")
	log.Println(r.Method, r.URL, http.StatusOK)
	w.Write(c.Status.ToJSON())
}

// databasesHandler returns the information of the running databases.
func (c *Conn) databasesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		c.createDatabaseHandler(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	log.Println(r.Method, r.URL, http.StatusOK)
	w.Write(c.Pila.Status().ToJSON())
}

// createDatabaseHandler creates a Database and returns 201 and the ID and name
// of the Database.
func (c *Conn) createDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		log.Println(r.Method, r.URL, http.StatusBadRequest, "missing name")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := pila.NewDatabase(name)
	err := c.Pila.AddDatabase(db)
	if err != nil {
		log.Println(r.Method, r.URL, http.StatusConflict, err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	log.Println(r.Method, r.URL, http.StatusCreated)
	w.WriteHeader(http.StatusCreated)
	w.Write(db.Status().ToJSON())
}

// databaseHandler returns the information of a single database given its ID
// or name.
func (c *Conn) databaseHandler(databaseID string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// we override the mux vars to be able to test
		// an arbitrary database ID
		if databaseID != "" {
			vars = map[string]string{
				"id": databaseID,
			}
		}

		db, ok := ResourceDatabase(c, vars["id"])
		if !ok {
			c.goneHandler(w, r, fmt.Sprintf("database %s is Gone", vars["id"]))
			return
		}

		if r.Method == "DELETE" {
			_ = c.Pila.RemoveDatabase(db.ID)
			log.Println(r.Method, r.URL, http.StatusNoContent)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		log.Println(r.Method, r.URL, http.StatusOK)
		w.Write(db.Status().ToJSON())
	})
}

// stacksHandler handles the stacks of a database, being able to get the status
// of them, or create a new one.
func (c *Conn) stacksHandler(databaseID string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// we override the mux vars to be able to test
		// an arbitrary database ID
		if databaseID != "" {
			vars = map[string]string{
				"database_id": databaseID,
			}
		}

		db, ok := ResourceDatabase(c, vars["database_id"])
		if !ok {
			c.goneHandler(w, r, fmt.Sprintf("database %s is Gone", vars["database_id"]))
			return
		}

		if r.Method == "PUT" {
			c.createStackHandler(w, r, db.ID.String())
			return
		}

		res, err := db.StacksStatus().ToJSON()
		if err != nil {
			log.Println(r.Method, r.URL, http.StatusBadRequest,
				"error on response serialization:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
		log.Println(r.Method, r.URL, http.StatusOK)

	})
}

// createStackHandler handles the creation of a stack, given a database
// by its id. Returns the status of the new stack.
func (c *Conn) createStackHandler(w http.ResponseWriter, r *http.Request, databaseID string) {
	name := r.FormValue("name")
	if name == "" {
		log.Println(r.Method, r.URL, http.StatusBadRequest, "missing name")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, ok := c.Pila.Database(uuid.UUID(databaseID))
	if !ok {
		c.goneHandler(w, r, fmt.Sprintf("database %s is Gone", databaseID))
		return
	}

	stack := pila.NewStack(name)
	err := db.AddStack(stack)
	if err != nil {
		log.Println(r.Method, r.URL, http.StatusConflict, err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	// Do not check error as the Status of a new stack does
	// not contain types that could cause such case.
	// See http://golang.org/src/encoding/json/encode.go?s=5438:5481#L125
	res, _ := stack.Status().ToJSON()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
	log.Println(r.Method, r.URL, http.StatusCreated)
}

// stackHandler handles operations on a single stack of a database. It holds
// the PUSH, POP and PEEK methods, and the stack deletion.
func (c *Conn) stackHandler(params *map[string]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		// we override the mux vars to be able to test
		// an arbitrary database and stack ID
		if params != nil {
			vars = *params
		}

		db, ok := ResourceDatabase(c, vars["database_id"])
		if !ok {
			c.goneHandler(w, r, fmt.Sprintf("database %s is Gone", vars["database_id"]))
			return
		}

		stack, ok := ResourceStack(db, vars["stack_id"])
		if !ok {
			c.goneHandler(w, r, fmt.Sprintf("stack %s is Gone", vars["stack_id"]))
			return
		}

		if r.Method == "POST" {
			c.pushStackHandler(w, r, stack)
			return
		}

		if r.Method == "DELETE" {
			_ = r.ParseForm()
			if _, ok := r.Form["flush"]; ok {
				c.flushStackHandler(w, r, stack)
				return
			}
			if _, ok := r.Form["full"]; ok {
				c.deleteStackHandler(w, r, db, stack)
				return
			}
			c.popStackHandler(w, r, stack)
			return
		}
	})
}

// pushStackHandler adds an element into a Stack and returns 200 and the element.
func (c *Conn) pushStackHandler(w http.ResponseWriter, r *http.Request, stack *pila.Stack) {
	if r.Body == nil {
		log.Println(r.Method, r.URL, http.StatusBadRequest,
			"no element provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var element pila.Element
	err := element.Decode(r.Body)
	if err != nil {
		log.Println(r.Method, r.URL, http.StatusBadRequest,
			"error on decoding element:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	stack.Push(element.Value)

	log.Println(r.Method, r.URL, http.StatusOK, element.Value)
	w.Header().Set("Content-Type", "application/json")

	// Do not check error as we consider our element
	// suitable for a JSON encoding.
	b, _ := element.ToJSON()
	w.Write(b)
}

// popStackHandler extracts the peek element of a Stack, returns 200 and returns it.
func (c *Conn) popStackHandler(w http.ResponseWriter, r *http.Request, stack *pila.Stack) {
	value, ok := stack.Pop()
	if !ok {
		log.Println(r.Method, r.URL, http.StatusNoContent)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	element := pila.Element{Value: value}

	log.Println(r.Method, r.URL, http.StatusOK, element.Value)
	w.Header().Set("Content-Type", "application/json")

	// Do not check error as we consider our element
	// suitable for a JSON encoding.
	b, _ := element.ToJSON()
	w.Write(b)
}

// flushStackHandler flushes the Stack, setting the size to 0 and emptying all
// the content.
func (c *Conn) flushStackHandler(w http.ResponseWriter, r *http.Request, stack *pila.Stack) {
	stack.Flush()

	log.Println(r.Method, r.URL, http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// Do not check error as we consider that a flushed
	// stack has no JSON encoding issues.
	b, _ := stack.Status().ToJSON()
	w.Write(b)
}

// deleteStackHandler deletes the Stack from a database.
func (c *Conn) deleteStackHandler(w http.ResponseWriter, r *http.Request, database *pila.Database, stack *pila.Stack) {
	stack.Flush()

	// Do not check output as we validated that
	// stack always exists.
	_ = database.RemoveStack(stack.ID)

	log.Println(r.Method, r.URL, http.StatusNoContent)
	w.WriteHeader(http.StatusNoContent)
	return
}

// notFoundHandler logs and returns a 404 NotFound response.
func (c *Conn) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL, http.StatusNotFound)
	http.NotFound(w, r)
}

// goneHandler logs and returns a 410 Gone response with information
// about the missing resource.
func (c *Conn) goneHandler(w http.ResponseWriter, r *http.Request, message string) {
	log.Println(r.Method, r.URL,
		http.StatusGone, message)
	w.WriteHeader(http.StatusGone)
}
