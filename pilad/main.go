// Binary pilad provides the daemon that runs the piladb server.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	flag.Parse()
	if versionFlag {
		fmt.Println(v())
		return
	}

	conn := NewConn()
	conn.buildConfig()
	logo(conn)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", conn.Config.Port()),
		Handler:      Router(conn),
		ReadTimeout:  conn.Config.ReadTimeout() * time.Second,
		WriteTimeout: conn.Config.WriteTimeout() * time.Second,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server stopped")
}
