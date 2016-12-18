package date

import (
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	inputOutput := []struct {
		input  time.Time
		output string
	}{
		{time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), "2009-11-10T23:00:00Z"},
		{time.Date(2016, time.May, 12, 15, 34, 56, 567464295, time.FixedZone("UTC", 3600)), "2016-05-12T15:34:56.567464295+01:00"},
	}

	for _, io := range inputOutput {
		if d := Format(io.input); d != io.output {
			t.Errorf("date is %s, expected %s", d, io.output)
		}
	}
}
