//go:build windows

package filesystem

import (
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

func creationTime(filename string) (time.Time, bool) {
	path, err := windows.UTF16PtrFromString(filename)
	if err != nil {
		return time.Time{}, false
	}

	var data windows.Win32FileAttributeData
	err = windows.GetFileAttributesEx(
		path,
		windows.GetFileExInfoStandard,
		(*byte)(unsafe.Pointer(&data)),
	)
	if err != nil {
		return time.Time{}, false
	}

	nanoseconds := data.CreationTime.Nanoseconds()
	if nanoseconds <= 0 {
		return time.Time{}, false
	}

	return time.Unix(0, nanoseconds), true
}
