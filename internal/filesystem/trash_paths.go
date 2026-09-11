package filesystem

import (
	"path"
	"path/filepath"
	"strings"
)

const (
	StateDirectoryName = ".serverigh"
	TrashDirectoryName = "trash"
	TrashPath          = "/trash"
)

func isTrashPath(requestPath string) bool {
	cleanPath := path.Clean("/" + strings.TrimPrefix(requestPath, "/"))

	return cleanPath == TrashPath ||
		strings.HasPrefix(cleanPath, TrashPath+"/")
}

func trashStoragePath(root string) string {
	return filepath.Join(root, StateDirectoryName, TrashDirectoryName)
}

func trashStorageChild(root string, requestPath string) string {
	relative := strings.TrimPrefix(requestPath, TrashPath)
	relative = strings.TrimPrefix(relative, "/")

	return filepath.Join(trashStoragePath(root), filepath.FromSlash(relative))
}
