package fileinfo

import "time"

const UnavailableTimeLabel = "Unavailable"

func FormatTimeValue(value time.Time) string {
	return value.Format(time.RFC3339)
}

func FormatTime(value time.Time) string {
	return value.Format("2006-01-02 15:04")
}

func FormatOptionalTime(value time.Time, ok bool) string {
	if !ok || value.IsZero() {
		return UnavailableTimeLabel
	}

	return FormatTime(value)
}

func FormatOptionalTimeValue(value time.Time, ok bool) string {
	if !ok || value.IsZero() {
		return ""
	}

	return FormatTimeValue(value)
}
