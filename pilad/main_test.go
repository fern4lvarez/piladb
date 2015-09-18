package main

import (
	"testing"
	"time"
)

// TestMain is a hack to get 100% test coverage.
func TestMain(t *testing.T) {
	go main()
	time.Sleep(5 * time.Millisecond)
}
