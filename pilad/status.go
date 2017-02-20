package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"
)

// Status represents the status of the running piladb
// instance.
type Status struct {
	Code             string    `json:"status"`
	Version          string    `json:"version"`
	GoVersion        string    `json:"go_version"`
	Host             string    `json:"host"`
	PID              int       `json:"pid"`
	StartedAt        time.Time `json:"started_at"`
	RunningFor       float64   `json:"running_for"`
	NumberGoroutines int       `json:"number_goroutines"`
	MemoryAlloc      string    `json:"memory_alloc"`
}

// NewStatus returns a new piladb status.
func NewStatus(version string, now time.Time, mem *runtime.MemStats) *Status {
	status := &Status{}
	status.Code = "OK"
	status.Version = version
	status.GoVersion = runtime.Version()
	status.Host = fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)
	status.PID = os.Getpid()
	status.StartedAt = now
	status.NumberGoroutines = runtime.NumGoroutine()
	if mem != nil {
		status.MemoryAlloc = MemOutput(mem.Alloc)
	}

	return status
}

// Update updates the Status given a current time and memory
// stats.
func (s *Status) Update(now time.Time, mem *runtime.MemStats) {
	s.RunningFor = now.Sub(s.StartedAt).Seconds()
	s.NumberGoroutines = runtime.NumGoroutine()
	s.MemoryAlloc = MemOutput(mem.Alloc)
}

// ToJSON returns the Status into a JSON file in []byte
// format.
func (s *Status) ToJSON() []byte {
	// Do not check error as the Status type does
	// not contain types that could cause such case.
	// See http://golang.org/src/encoding/json/encode.go?s=5438:5481#L125
	s.StartedAt = s.StartedAt.Local()
	b, _ := json.Marshal(s)
	return b
}
