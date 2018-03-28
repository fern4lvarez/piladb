package main

import (
	"os"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	os.Setenv("PILADB_PORT", "35343")
	go main()
	time.Sleep(5 * time.Millisecond)
	os.Setenv("PILADB_PORT", "")
}

func TestMain_Error(t *testing.T) {
	os.Setenv("PILADB_PORT", "35343")
	go main()
	go main()
	time.Sleep(5 * time.Millisecond)
	os.Setenv("PILADB_PORT", "")
}

func TestMain_Version(t *testing.T) {
	versionFlag = true
	go main()
	t.Log(v())
}
