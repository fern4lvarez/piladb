package main

import (
	"os"
	"testing"
)

func TestPort(t *testing.T) {
	if port := Port(); port != "1205" {
		t.Errorf("port is %s, expected %s", port, "1205")
	}
}

func TestPort_Env(t *testing.T) {
	os.Setenv("PILADB_PORT", "8888")
	if port := Port(); port != "8888" {
		t.Errorf("port is %s, expected %s", port, "8888")
	}
}
