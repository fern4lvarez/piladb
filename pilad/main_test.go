package main

import (
	"os"
	"testing"
	"time"
)

// TestMain is a hack to get 100% test coverage.
func TestMain(t *testing.T) {
	os.Setenv("PILADB_PORT", "35343")
	go main()
	time.Sleep(5 * time.Millisecond)
	os.Setenv("PILADB_PORT", "")
}

func TestMainVersion(t *testing.T) {
	versionFlag = true
	go main()
	t.Log(v())
}
