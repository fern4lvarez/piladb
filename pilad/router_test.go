package main

import "testing"

func TestRouter(t *testing.T) {
	conn := NewConn()
	if r := Router(conn); r == nil {
		t.Fatal("router is nil")
	}

}
