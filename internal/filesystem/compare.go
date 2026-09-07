package filesystem

import (
	"strings"
	"time"
)

func compareNames(left string, right string) int {
	return strings.Compare(strings.ToLower(left), strings.ToLower(right))
}

func compareInt64(left int64, right int64) int {
	switch {
	case left < right:
		return -1
	case left > right:
		return 1
	default:
		return 0
	}
}

func compareTimes(left time.Time, right time.Time) int {
	switch {
	case left.Before(right):
		return -1
	case left.After(right):
		return 1
	default:
		return 0
	}
}

func compareOptionalTimes(
	left time.Time,
	leftKnown bool,
	right time.Time,
	rightKnown bool,
) int {
	if leftKnown != rightKnown {
		if leftKnown {
			return -1
		}

		return 1
	}

	return compareTimes(left, right)
}
