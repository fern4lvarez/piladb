package date

import "time"

// Format gets a time.Time and formats it to RFC3339, which is the default format for piladb responses.
func Format(t time.Time) string {
	return t.Format("2006-01-02T15:04:05.999999999-07:00")
}
