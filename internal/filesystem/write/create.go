package write

import (
	"errors"
	"os"
	"path/filepath"
)

func (s Service) CreateDirectory(parentPath string, name string) error {
	target, err := s.directoryTarget(parentPath, name)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0o755); err != nil {
		return wrapWriteError("create directory", err)
	}

	return nil
}

func (s Service) CreateFile(parentPath string, name string) error {
	target, err := s.childTarget(parentPath, name)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return wrapWriteError("create file", err)
	}

	if err := file.Close(); err != nil {
		return wrapWriteError("close created file", err)
	}

	return nil
}

func (s Service) directoryTarget(parentPath string, name string) (string, error) {
	relativePath, err := cleanDirectoryPath(name)
	if err != nil {
		return "", err
	}

	parent, err := s.resolve(parentPath)
	if err != nil {
		return "", err
	}
	if !parent.Info.IsDir() {
		return "", ErrNotDirectory
	}

	target := filepath.Join(parent.Absolute, filepath.FromSlash(relativePath))
	if err := s.ensureInsideRoot(target); err != nil {
		return "", err
	}
	if _, err := os.Lstat(target); err == nil {
		return "", ErrPathExists
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", wrapFileError("inspect target path", err)
	}

	return target, nil
}

func (s Service) Rename(requestPath string, name string) error {
	source, err := s.resolve(requestPath)
	if err != nil {
		return err
	}

	target, err := s.childTargetForAbsolute(filepath.Dir(source.Absolute), name)
	if err != nil {
		return err
	}

	if err := os.Rename(source.Absolute, target); err != nil {
		return wrapWriteError("rename path", err)
	}

	return nil
}
