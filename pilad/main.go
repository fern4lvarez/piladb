// Binary pilad provides the daemon that runs the piladb server.
package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	flag.Parse()
	if versionFlag {
		fmt.Println(v())
		return
	}

	if err := start(NewConn()); err != nil {
		log.Println(err)
	}
}
