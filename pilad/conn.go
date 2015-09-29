package main

import (
	"log"
	"net/http"
	"time"

	"github.com/fern4lvarez/piladb/pila"
	"github.com/fern4lvarez/piladb/pkg/version"
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
