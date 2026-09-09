package write

import (
	"os"
	"path/filepath"
)

func (s Service) CreateDirectory(parentPath string, name string) error {
	target, err := s.childTarget(parentPath, name)
	if err != nil {
		return err
	}

	if err := os.Mkdir(target, 0o755); err != nil {
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
