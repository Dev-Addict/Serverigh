//go:build windows

package fileinfo

import (
	"strings"

	"golang.org/x/sys/windows"
)

func IsHidden(filename string, name string) bool {
	if strings.HasPrefix(name, ".") {
		return true
	}

	path, err := windows.UTF16PtrFromString(filename)
	if err != nil {
		return false
	}

	attributes, err := windows.GetFileAttributes(path)
	if err != nil {
		return false
	}

	return attributes&windows.FILE_ATTRIBUTE_HIDDEN != 0
}
