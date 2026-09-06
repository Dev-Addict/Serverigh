package filesystem

import "time"

const unavailableTimeLabel = "Unavailable"

func formatTime(value time.Time) string {
	return value.Format("2006-01-02 15:04")
}

func formatOptionalTime(value time.Time, ok bool) string {
	if !ok || value.IsZero() {
		return unavailableTimeLabel
	}

	return formatTime(value)
}
