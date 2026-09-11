package write

import (
	"path"
	"strings"
)

func cleanDirectoryPath(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" || path.IsAbs(name) || strings.ContainsAny(name, "\\\x00") {
		return "", ErrInvalidPath
	}

	cleanPath := path.Clean(name)
	if cleanPath == "." || cleanPath == ".." ||
		strings.HasPrefix(cleanPath, "../") {
		return "", ErrInvalidPath
	}
	for _, segment := range strings.Split(cleanPath, "/") {
		if err := validateEntryName(segment); err != nil {
			return "", err
		}
	}

	return cleanPath, nil
}
