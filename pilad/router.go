package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Router returns a gorila/mux Router with all specified endpoints and
// handlers.
func Router(conn *Conn) *mux.Router {
	r := mux.NewRouter()
	conn.opDate = time.Now()

	// GET /
	r.HandleFunc("/", conn.rootHandler).
		Methods("GET")

	// GET /_status
	r.HandleFunc("/_status", conn.statusHandler).
		Methods("GET")

	// GET /_config
	r.HandleFunc("/_config", conn.configHandler).
		Methods("GET")
	// GET /_config/$CONFIG_KEY
	// POST /_config/$CONFIG_KEY + {element: value}
	r.Handle("/_config/{key}", conn.configKeyHandler("")).
		Methods("GET", "POST")

	// GET /databases
	// PUT /databases?name=DATABASE_NAME
	r.HandleFunc("/databases", conn.databasesHandler).
		Methods("GET", "PUT")
	// GET /databases/$DATABASE_ID
	// DELETE /databases/$DATABASE_ID
	r.Handle("/databases/{id}", conn.databaseHandler("")).
		Methods("GET", "DELETE")

	// GET /databases/$DATABASE_ID/stacks
	// GET /databases/$DATABASE_ID/stacks?kv
	// PUT /databases/$DATABASE_ID/stacks?name=STACK_NAME
	r.Handle("/databases/{database_id}/stacks", conn.stacksHandler("")).
		Methods("GET", "PUT")

	// GET /databases/$DATABASE_ID/stacks/$STACK_ID
	// GET /databases/$DATABASE_ID/stacks/$STACK_ID?peek
	// GET /databases/$DATABASE_ID/stacks/$STACK_ID?size
	// POST /databases/$DATABASE_ID/stacks/$STACK_ID + {element: value}
	// DELETE /databases/$DATABASE_ID/stacks/$STACK_ID
	// DELETE /databases/$DATABASE_ID/stacks/$STACK_ID?flush
	// DELETE /databases/$DATABASE_ID/stacks/$STACK_ID?full
	r.Handle("/databases/{database_id}/stacks/{stack_id}", conn.stackHandler(nil)).
		Methods("GET", "POST", "DELETE")

	r.NotFoundHandler = http.HandlerFunc(conn.notFoundHandler)
	return r
}
