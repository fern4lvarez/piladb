package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fern4lvarez/piladb/vendor/src/github.com/gorilla/mux"
)

func main() {
	conn := NewConn()
	r := mux.NewRouter()
	r.HandleFunc("/_status", conn.statusHandler).
		Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(conn.notFoundHandler)

	log.Printf("piladb is listening to port %s", Port())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", Port()), r))
}
