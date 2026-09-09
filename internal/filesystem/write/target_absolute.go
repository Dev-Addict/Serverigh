package write

import (
	"errors"
	"os"
	"path/filepath"
)

func (s Service) childTargetForAbsolute(parent string, name string) (string, error) {
	if err := validateEntryName(name); err != nil {
		return "", err
	}

	target := filepath.Join(parent, name)
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
