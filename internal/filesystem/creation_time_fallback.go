//go:build !linux && !darwin && !freebsd && !netbsd && !windows

package filesystem

import "time"

func creationTime(filename string) (time.Time, bool) {
	return time.Time{}, false
}
