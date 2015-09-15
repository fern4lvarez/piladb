package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	conn := NewConn()
	http.HandleFunc("/_status", conn.statusHandler)

	log.Printf("piladb is listening to port %s", Port())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", Port()), nil))
}
