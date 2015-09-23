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
	w.Header().Set("Content-Type", "application/json")
	log.Println(r.Method, r.URL, http.StatusOK)
	w.Write(c.Pila.Status())
}

// notFoundHandler logs and returns a 404 NotFound response.
func (c *Conn) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL, http.StatusNotFound)
	http.NotFound(w, r)
}
