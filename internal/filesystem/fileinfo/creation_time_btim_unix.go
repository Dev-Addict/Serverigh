//go:build darwin || freebsd || netbsd

package fileinfo

import (
	"time"

	"golang.org/x/sys/unix"
)

func CreationTime(filename string) (time.Time, bool) {
	var stat unix.Stat_t
	if err := unix.Stat(filename, &stat); err != nil {
		return time.Time{}, false
	}

	if stat.Btim.Sec <= 0 {
		return time.Time{}, false
	}

	return time.Unix(int64(stat.Btim.Sec), int64(stat.Btim.Nsec)), true
}
