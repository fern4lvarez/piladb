package main

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/fern4lvarez/piladb/pkg/date"
)

func TestNewStatus(t *testing.T) {
	now := time.Now()
	mem := runtime.MemStats{Alloc: 0}
	status := NewStatus("v1", now, &mem)

	if status == nil {
		t.Fatal("status is nil")
	}
	if status.Code != "OK" {
		t.Errorf("status is %s, expected %s", status.Code, "OK")
	}
	if status.Version != "v1" {
		t.Errorf("version is %s, expected %s", status.Version, "v1")
	}
	if host := fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH); status.Host != host {
		t.Errorf("host is %s, expected %s", status.Host, host)
	}
	if status.PID != os.Getpid() {
		t.Errorf("PID is %d, expected %d", status.PID, os.Getpid())
	}
	if status.StartedAt != now {
		t.Errorf("version is %v expected %v", status.StartedAt, now)
	}
	if status.NumberGoroutines != runtime.NumGoroutine() {
		t.Errorf("number of goroutines is %v expected %v", status.StartedAt, now)
	}
	if status.MemoryAlloc != "0B" {
		t.Errorf("memory allocated is %v, expected to be %v", status.MemoryAlloc, "0.00MB")
	}
}

func TestStatusUpdate(t *testing.T) {
	now := time.Now()
	oneHourAgo := now.Add(-60 * time.Minute)
	status := NewStatus("v1", oneHourAgo, nil)

	mem := runtime.MemStats{Alloc: 7353735469}
	numberGoroutines := runtime.NumGoroutine()
	status.Update(now, &mem)

	if r := status.RunningFor; r != 3600.0 {
		t.Errorf("running for is %v, expected %v", r, 3600.0)
	}
	if n := status.NumberGoroutines; n != numberGoroutines {
		t.Errorf("number of goroutines is %v, expected %v", n, numberGoroutines)
	}
	if m := status.MemoryAlloc; m != "6.85GiB" {
		t.Errorf("memory allocated is %v, expected %v", m, "6.85GiB")
	}
}

func TestStatusToJSON(t *testing.T) {
	now := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	status := NewStatus("v1", now, nil)
	oneHourLater := now.Add(60 * time.Minute)
	mem := runtime.MemStats{Alloc: 0}
	expectedJSON := fmt.Sprintf(`{"status":"OK","version":"v1","go_version":"%s","host":"%s_%s","pid":%d,"started_at":"%s","running_for":3600,"number_goroutines":%d,"memory_alloc":"0B"}`, runtime.Version(), runtime.GOOS, runtime.GOARCH, os.Getpid(), date.Format(now.Local()), runtime.NumGoroutine())

	status.Update(oneHourLater, &mem)
	json := status.ToJSON()
	if json == nil {
		t.Fatal("json is nil")
	}
	if string(json) != expectedJSON {
		t.Errorf("json is %s, expected %s", string(json), expectedJSON)
	}
}
