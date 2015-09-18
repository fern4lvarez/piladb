package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	conn := NewConn()
	log.Printf("piladb is listening to port %s", Port())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", Port()), Router(conn)))
}
