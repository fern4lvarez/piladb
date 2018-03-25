package main

import (
	"net/http"
	"syscall"
	"testing"
)

func TestStart(t *testing.T) {
	conn := NewConn()
	conn.stop <- syscall.SIGTERM
	if err := start(conn); err != nil {
		t.Errorf("err is %v, expected nil", err)
	}
}

func TestListenGracefulShutdown(t *testing.T) {
	conn := NewConn()
	conn.srv = &http.Server{}

	go listenGracefulShutdown(conn)
	conn.stop <- syscall.SIGTERM
}
