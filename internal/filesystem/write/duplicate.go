package write

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (s Service) Duplicate(requestPath string) error {
	source, err := s.resolve(requestPath)
	if err != nil {
		return err
	}

	target, err := s.uniqueDuplicateTarget(
		filepath.Dir(source.Absolute),
		source.Info.Name(),
	)
	if err != nil {
		return err
	}

	if source.Info.IsDir() {
		err = copyDirectory(source.Absolute, target)
	} else {
		err = copyFile(source.Absolute, target, source.Info.Mode())
	}
	if err != nil {
		return wrapWriteError("duplicate path", err)
	}

	return nil
}

func (s Service) uniqueDuplicateTarget(parent string, name string) (string, error) {
	if err := validateEntryName(name); err != nil {
		return "", err
	}

	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	for index := 1; index < 10_000; index++ {
		target := filepath.Join(parent, duplicateName(base, ext, index))
		if err := s.ensureInsideRoot(target); err != nil {
			return "", err
		}
		if _, err := os.Lstat(target); os.IsNotExist(err) {
			return target, nil
		} else if err != nil {
			return "", wrapFileError("inspect duplicate path", err)
		}
	}

	return "", ErrPathExists
}

func duplicateName(base string, ext string, index int) string {
	return base + " " + strconv.Itoa(index) + ext
}
