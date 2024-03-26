package helpers

import "time"

func Now() *time.Time {
	current, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	return &current
}
