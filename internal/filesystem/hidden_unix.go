//go:build !windows

package filesystem

import "strings"

func isHidden(_ string, name string) bool {
	return strings.HasPrefix(name, ".")
}
