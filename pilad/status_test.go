package main

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestNewStatus(t *testing.T) {
	now := time.Now()
	status := NewStatus("v1", now)

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
}

func TestStatusSetRunningFor(t *testing.T) {
	now := time.Now()
	oneHourAgo := now.Add(-60 * time.Minute)
	status := NewStatus("v1", oneHourAgo)

	if r := status.SetRunningFor(now); r != 3600.0 {
		t.Errorf("running for is %v, expected %v", r, 3600.0)
	}
	if r := status.RunningFor; r != 3600.0 {
		t.Errorf("running for is %v, expected %v", r, 3600.0)
	}

}

func TestStatusToJSON(t *testing.T) {
	now := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	status := NewStatus("v1", now)
	oneHourLater := now.Add(60 * time.Minute)
	expectedJSON := fmt.Sprintf(`{"status":"OK","version":"v1","host":"%s_%s","pid":%d,"started_at":"2009-11-10T23:00:00Z","running_for":3600}`, runtime.GOOS, runtime.GOARCH, os.Getpid())

	json := status.ToJSON(oneHourLater)
	if json == nil {
		t.Fatal("json is nil")
	}
	if string(json) != expectedJSON {
		t.Errorf("json is %s, expected %s", string(json), expectedJSON)
	}
}
