package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	conn := NewConn()
	r := Router(conn)
	log.Printf("piladb is listening to port %s", Port())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", Port()), r))
}
