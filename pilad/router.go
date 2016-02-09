package main

import (
	"net/http"

	"github.com/fern4lvarez/piladb/_vendor/src/github.com/gorilla/mux"
)

// Router returns a gorila/mux Router with all specified endpoints and
// handlers.
func Router(conn *Conn) *mux.Router {
	r := mux.NewRouter()

	// GET /_status
	r.HandleFunc("/_status", conn.statusHandler).
		Methods("GET")

	// GET /databases
	// PUT /databases?name=DATABASE_NAME
	r.HandleFunc("/databases", conn.databasesHandler).
		Methods("GET", "PUT")
	// GET /databases/$DATABASE_ID
	r.Handle("/databases/{id}", conn.databaseHandler("")).
		Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(conn.notFoundHandler)
	return r
}
