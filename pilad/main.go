// Binary pilad provides the daemon that runs the piladb server.
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	conn := NewConn()
	logo(conn)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", Port()), Router(conn)))
}
