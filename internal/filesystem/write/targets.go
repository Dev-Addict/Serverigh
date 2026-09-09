package write

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (s Service) childTarget(parentPath string, name string) (string, error) {
	if err := validateEntryName(name); err != nil {
		return "", err
	}

	parent, err := s.resolve(parentPath)
	if err != nil {
		return "", err
	}
	if !parent.Info.IsDir() {
		return "", ErrNotDirectory
	}

	return s.childTargetForAbsolute(parent.Absolute, name)
}

func (s Service) uniqueChildTarget(parent string, name string) (string, error) {
	if err := validateEntryName(name); err != nil {
		return "", err
	}

	target := filepath.Join(parent, name)
	if err := s.ensureInsideRoot(target); err != nil {
		return "", err
	}
	if _, err := os.Lstat(target); errors.Is(err, os.ErrNotExist) {
		return target, nil
	} else if err != nil {
		return "", wrapFileError("inspect target path", err)
	}

	return uniqueCopyTarget(parent, name)
}

func uniqueCopyTarget(parent string, name string) (string, error) {
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	for index := 1; index < 10_000; index++ {
		target := filepath.Join(parent, copyName(base, ext, index))
		if _, err := os.Lstat(target); errors.Is(err, os.ErrNotExist) {
			return target, nil
		} else if err != nil {
			return "", wrapFileError("inspect target path", err)
		}
	}

	return "", ErrPathExists
}

func copyName(base string, ext string, index int) string {
	if index == 1 {
		return base + " copy" + ext
	}

	return base + " copy " + strconv.Itoa(index) + ext
}

func sameOrDescendant(candidate string, parent string) bool {
	relative, err := filepath.Rel(parent, candidate)
	if err != nil {
		return false
	}

	return relative == "." ||
		(relative != ".." &&
			!strings.HasPrefix(relative, ".."+string(filepath.Separator)))
}
