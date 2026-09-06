//go:build linux

package filesystem

import (
	"time"

	"golang.org/x/sys/unix"
)

func creationTime(filename string) (time.Time, bool) {
	var stat unix.Statx_t
	err := unix.Statx(
		unix.AT_FDCWD,
		filename,
		0,
		unix.STATX_BTIME,
		&stat,
	)
	if err != nil {
		return time.Time{}, false
	}

	if stat.Mask&unix.STATX_BTIME == 0 {
		return time.Time{}, false
	}

	if stat.Btime.Sec <= 0 {
		return time.Time{}, false
	}

	return time.Unix(stat.Btime.Sec, int64(stat.Btime.Nsec)), true
}
