package main

import (
	"net/http"

	"github.com/gorilla/mux"
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
	// DELETE /databases/$DATABASE_ID
	r.Handle("/databases/{id}", conn.databaseHandler("")).
		Methods("GET", "DELETE")

	// GET /databases/$DATABASE_ID/stacks
	// PUT /databases/$DATABASE_ID/stacks?name=STACK_NAME
	r.Handle("/databases/{database_id}/stacks", conn.stacksHandler("")).
		Methods("GET", "PUT")

	// GET /databases/$DATABASE_ID/stacks/$STACK_ID
	// GET /databases/$DATABASE_ID/stacks/$STACK_ID?full
	// POST /databases/$DATABASE_ID/stacks/$STACK_ID + {element: value}
	// DELETE /databases/$DATABASE_ID/stacks/$STACK_ID
	// DELETE /databases/$DATABASE_ID/stacks/$STACK_ID?flush
	// DELETE /databases/$DATABASE_ID/stacks/$STACK_ID?full
	r.Handle("/databases/{database_id}/stacks/{stack_id}", conn.stackHandler(nil)).
		Methods("GET", "POST", "DELETE")

	r.NotFoundHandler = http.HandlerFunc(conn.notFoundHandler)
	return r
}
