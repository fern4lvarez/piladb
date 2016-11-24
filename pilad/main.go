// Binary pilad provides the daemon that runs the piladb server.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	flag.Parse()
	conn := NewConn()
	conn.buildConfig()
	logo(conn)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", Port()),
		Handler:      Router(conn),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
