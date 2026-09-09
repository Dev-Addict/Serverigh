package write

import "strings"

func validateEntryName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" || name == "." || name == ".." {
		return ErrInvalidPath
	}
	if strings.ContainsAny(name, `/\`) || strings.ContainsRune(name, '\x00') {
		return ErrInvalidPath
	}

	return nil
}
