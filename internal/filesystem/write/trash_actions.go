package write

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

func (s Service) RestoreTrash(requestPath string) error {
	source, err := s.resolveTrashEntry(requestPath)
	if err != nil {
		return err
	}

	metadata, err := s.readTrashMetadata(path.Base(source.Path))
	if err != nil {
		return err
	}

	destination, err := s.restoreDestination(metadata.OriginalPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return wrapWriteError("create restore parent directory", err)
	}
	if err := os.Rename(source.Absolute, destination); err != nil {
		return wrapWriteError("restore trash entry", err)
	}

	return s.removeTrashMetadata(path.Base(source.Path))
}

func (s Service) DeleteTrash(requestPath string) error {
	source, err := s.resolveTrashEntry(requestPath)
	if err != nil {
		return err
	}
	if err := os.RemoveAll(source.Absolute); err != nil {
		return wrapWriteError("delete trash entry", err)
	}

	return s.removeTrashMetadata(path.Base(source.Path))
}

func (s Service) resolveTrashEntry(requestPath string) (ResolvedPath, error) {
	source, err := s.resolve(requestPath)
	if err != nil {
		return ResolvedPath{}, err
	}
	if source.Path == trashPath || !strings.HasPrefix(source.Path, trashPath+"/") {
		return ResolvedPath{}, ErrInvalidPath
	}
	if path.Base(source.Path) == trashMetaFolderName {
		return ResolvedPath{}, ErrInvalidPath
	}

	return source, nil
}

func (s Service) restoreDestination(originalPath string) (string, error) {
	destinationPath, err := cleanRequestPath(originalPath)
	if err != nil {
		return "", err
	}

	destination := filepath.Join(
		s.root,
		filepath.FromSlash(path.Clean(strings.TrimPrefix(destinationPath, "/"))),
	)
	if err := s.ensureInsideRoot(destination); err != nil {
		return "", err
	}
	if _, err := os.Lstat(destination); err == nil {
		return "", ErrPathExists
	} else if err != nil && !os.IsNotExist(err) {
		return "", wrapWriteError("inspect restore destination", err)
	}

	return destination, nil
}
