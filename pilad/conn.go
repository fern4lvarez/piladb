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

	"github.com/fern4lvarez/piladb/_vendor/src/github.com/gorilla/mux"
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
	w.Write(c.Status.ToJson(time.Now()))
}

// databasesHandler returns the information of the running databases.
func (c *Conn) databasesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		c.createDatabaseHandler(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	log.Println(r.Method, r.URL, http.StatusOK)
	w.Write(c.Pila.Status())
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
	w.Write(db.Status())
}

// notFoundHandler logs and returns a 404 NotFound response.
func (c *Conn) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL, http.StatusNotFound)
	http.NotFound(w, r)
}

// databaseHandler returns the information of a single database.
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
			log.Println(r.Method, r.URL,
				http.StatusGone, "database is Gone")
			w.WriteHeader(http.StatusGone)
			w.Write([]byte(fmt.Sprintf("database %s is Gone", vars["id"])))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		log.Println(r.Method, r.URL, http.StatusOK)
		w.Write(db.Status())
	})
}

// popStackHandler returns 200 and the first element of a Stack.
func (c *Conn) popStackHandler(params map[string]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		dbID := uuid.UUID(params["database_id"])
		db, ok := c.Pila.Database(dbID)
		if !ok {
			log.Println(r.Method, r.URL,
				http.StatusGone, "database is Gone")
			w.WriteHeader(http.StatusGone)
			w.Write([]byte(fmt.Sprintf("database %s is Gone", params["database_id"])))
			return
		}

		stackID := uuid.UUID(params["stack_id"])
		stack, ok := db.Stacks[stackID]
		if !ok {
			log.Println(r.Method, r.URL,
				http.StatusGone, "stack is Gone")
			w.WriteHeader(http.StatusGone)
			w.Write([]byte(fmt.Sprintf("stack %s is Gone", params["stack_id"])))
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
