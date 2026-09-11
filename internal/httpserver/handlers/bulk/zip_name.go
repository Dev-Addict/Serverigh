package bulk

import (
	"path"
	"strings"
)

func archiveName(requestPath string) string {
	name := strings.TrimPrefix(path.Clean(requestPath), "/")
	if name == "" || name == "." {
		return "root"
	}

	return name
}
