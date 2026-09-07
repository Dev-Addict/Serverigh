//go:build !linux && !darwin && !freebsd && !netbsd && !windows

package fileinfo

import "time"

func CreationTime(filename string) (time.Time, bool) {
	return time.Time{}, false
}
