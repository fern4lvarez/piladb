package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func start(conn *Conn) error {
	conn.buildConfig()
	conn.srv = &http.Server{
		Addr:         fmt.Sprintf(":%d", conn.Config.Port()),
		Handler:      Router(conn),
		ReadTimeout:  conn.Config.ReadTimeout() * time.Second,
		WriteTimeout: conn.Config.WriteTimeout() * time.Second,
	}
	logo(conn)

	go listenGracefulShutdown(conn)

	if err := conn.srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	<-conn.idle
	return nil
}

func listenGracefulShutdown(conn *Conn) {
	signal.Notify(conn.stop, os.Interrupt, syscall.SIGTERM)
	<-conn.stop

	log.Printf("Shutting down pilad and draining connections with a timeout of %s...", conn.srv.ReadTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), conn.srv.ReadTimeout)
	defer cancel()

	conn.srv.Shutdown(ctx)
	close(conn.idle)
	log.Println("Bye!")
}
