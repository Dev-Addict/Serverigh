package write

import (
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	stateFolderName = ".serverigh"
	trashFolderName = "trash"
)

func (s Service) Trash(requestPath string) error {
	source, err := s.resolve(requestPath)
	if err != nil {
		return err
	}

	trashPath := filepath.Join(s.root, stateFolderName, trashFolderName)
	if sameOrDescendant(trashPath, source.Absolute) {
		return ErrInvalidPath
	}

	if err := os.MkdirAll(trashPath, 0o755); err != nil {
		return wrapWriteError("create trash directory", err)
	}

	target := uniqueTrashPath(trashPath, path.Base(source.Path))
	if err := os.Rename(source.Absolute, target); err != nil {
		return wrapWriteError("move path to trash", err)
	}

	return nil
}

func uniqueTrashPath(trashPath string, name string) string {
	stamp := time.Now().UTC().Format("20060102T150405.000000000Z")

	return filepath.Join(trashPath, stamp+"-"+name)
}
