package main

import (
	"testing"
	"time"
)

func TestNewStatus(t *testing.T) {
	now := time.Now()
	status := NewStatus(now)

	if status == nil {
		t.Fatal("status is nil")
	}
	if status.Code != "OK" {
		t.Errorf("status is %s, expected %s", status.Code, "OK")
	}
	if status.Version != "0" {
		t.Errorf("version is %s, expected %s", status.Version, "0")
	}
	if status.StartedAt != now {
		t.Errorf("version is %v expected %v", status.StartedAt, now)
	}
}

func TestStatusSetRunningFor(t *testing.T) {
	now := time.Now()
	oneHourAgo := now.Add(-60 * time.Minute)
	status := NewStatus(oneHourAgo)

	if r := status.SetRunningFor(now); r != 3600.0 {
		t.Errorf("running for is %v, expected %v", r, 3600.0)
	}
	if r := status.RunningFor; r != 3600.0 {
		t.Errorf("running for is %v, expected %v", r, 3600.0)
	}

}

func TestStatusToJson(t *testing.T) {
	now := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	status := NewStatus(now)
	oneHourLater := now.Add(60 * time.Minute)
	expectedJSON := `{"status":"OK","version":"0","started_at":"2009-11-10T23:00:00Z","running_for":3600}`

	json := status.ToJson(oneHourLater)
	if json == nil {
		t.Fatal("json is nil")
	}
	if string(json) != expectedJSON {
		t.Errorf("json is %s, expected %s", string(json), expectedJSON)
	}
}
