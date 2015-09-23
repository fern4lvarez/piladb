package main

import (
	"net/http"

	"github.com/fern4lvarez/piladb/_vendor/src/github.com/gorilla/mux"
)

// Router returns a gorila/mux Router with all specified endpoints and
// handlers.
func Router(conn *Conn) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/_status", conn.statusHandler).
		Methods("GET")
	r.HandleFunc("/databases", conn.databasesHandler).
		Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(conn.notFoundHandler)
	return r
}
