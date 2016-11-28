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
		Addr:         fmt.Sprintf(":%d", conn.Config.Port()),
		Handler:      Router(conn),
		ReadTimeout:  conn.Config.ReadTimeout() * time.Second,
		WriteTimeout: conn.Config.WriteTimeout() * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
