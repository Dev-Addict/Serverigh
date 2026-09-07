//go:build !windows

package fileinfo

import "strings"

func IsHidden(_ string, name string) bool {
	return strings.HasPrefix(name, ".")
}
