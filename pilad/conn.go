package main

import (
	"encoding/json"
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
	conn.Status = NewStatus(version.CommitHash(), time.Now())
	return conn
}

// Connection Handlers

// statusHandler writes the piladb status into the response.
func (c *Conn) statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Println(r.Method, r.URL, http.StatusOK)
	w.Write(c.Status.ToJSON(time.Now()))
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
		log.Println(r.Method, r.URL, http.StatusConflict, "database exists")
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

		db, ok := c.Pila.Database(uuid.UUID(vars["id"]))
		if !ok {
			// Fallback to find by database name
			db, ok = c.Pila.Database(uuid.New(vars["id"]))
		}
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

		db, ok := c.Pila.Database(uuid.UUID(vars["database_id"]))
		if !ok {
			// Fallback to find by database name
			db, ok = c.Pila.Database(uuid.New(vars["database_id"]))
		}
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
				"error on response serialization")
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
		log.Println(r.Method, r.URL, http.StatusConflict, fmt.Sprintf("stack %s exists in database %s", name, databaseID))
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

// popStackHandler returns 200 and the first element of a Stack.
func (c *Conn) popStackHandler(params map[string]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		dbID := uuid.UUID(params["database_id"])
		db, ok := c.Pila.Database(dbID)
		if !ok {
			c.goneHandler(w, r, fmt.Sprintf("database %s is Gone", params["database_id"]))
			return
		}

		stackID := uuid.UUID(params["stack_id"])
		stack, ok := db.Stacks[stackID]
		if !ok {
			c.goneHandler(w, r, fmt.Sprintf("stack %s is Gone", params["stack_id"]))
			return
		}
		element, ok := stack.Pop()
		if !ok {
			log.Println(r.Method, r.URL, http.StatusNoContent)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		log.Println(r.Method, r.URL, http.StatusOK)

		b, _ := json.Marshal(element)
		w.Write(b)
	})
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
