package write

import (
	"path"
	"path/filepath"
	"strings"
)

const (
	stateFolderName     = ".serverigh"
	trashFolderName     = "trash"
	trashPath           = "/trash"
	trashMetaFolderName = "trash-meta"
)

func isTrashPath(requestPath string) bool {
	cleanPath := path.Clean("/" + strings.TrimPrefix(requestPath, "/"))

	return cleanPath == trashPath ||
		strings.HasPrefix(cleanPath, trashPath+"/")
}

func trashStoragePath(root string) string {
	return filepath.Join(root, stateFolderName, trashFolderName)
}

func trashMetaStoragePath(root string) string {
	return filepath.Join(root, stateFolderName, trashMetaFolderName)
}

func trashStorageChild(root string, requestPath string) string {
	relative := strings.TrimPrefix(requestPath, trashPath)
	relative = strings.TrimPrefix(relative, "/")

	return filepath.Join(trashStoragePath(root), filepath.FromSlash(relative))
}

func trashMetaPath(root string, name string) string {
	return filepath.Join(trashMetaStoragePath(root), name+".json")
}
