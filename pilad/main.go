// Binary pilad provides the daemon that runs the piladb server.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	flag.Parse()
	conn := NewConn()
	conn.buildConfig()
	logo(conn)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", Port()), Router(conn)))
}
