package searcher

import (
	"path"
	"strings"
)

func cleanRequestPath(requestPath string) (string, error) {
	if strings.ContainsRune(requestPath, '\x00') {
		return "", ErrInvalidPath
	}

	if requestPath == "" {
		requestPath = "/"
	}

	slashedPath := strings.ReplaceAll(requestPath, "\\", "/")
	for _, segment := range strings.Split(slashedPath, "/") {
		if segment == ".." {
			return "", ErrOutsideRoot
		}
	}

	cleanPath := path.Clean("/" + strings.TrimPrefix(slashedPath, "/"))
	if cleanPath == "." {
		return "/", nil
	}

	return cleanPath, nil
}
