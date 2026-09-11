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
	return s.uniqueNumberedTarget(parent, name, nil)
}

func (s Service) uniqueChildTargetForAbsolute(parent string, name string) (string, error) {
	target, err := s.childTargetForAbsolute(parent, name)
	if err == nil {
		return target, nil
	}
	if err != ErrPathExists {
		return "", err
	}

	return s.uniqueNumberedTarget(parent, name, nil)
}

func (s Service) uniqueChildNameForAbsolute(
	parent string,
	name string,
	reserved map[string]bool,
) (string, error) {
	if err := validateEntryName(name); err != nil {
		return "", err
	}
	if reserved == nil || !reserved[name] {
		target, err := s.childTargetForAbsolute(parent, name)
		if err == nil {
			return filepath.Base(target), nil
		}
		if err != ErrPathExists {
			return "", err
		}
	}

	target, err := s.uniqueNumberedTarget(parent, name, reserved)
	if err != nil {
		return "", err
	}

	return filepath.Base(target), nil
}

func (s Service) uniqueNumberedTarget(
	parent string,
	name string,
	reserved map[string]bool,
) (string, error) {
	if err := validateEntryName(name); err != nil {
		return "", err
	}

	base, ext, start := numberedNameParts(name)
	for index := start; index < 10_000; index++ {
		candidate := numberedName(base, ext, index)
		if reserved != nil && reserved[candidate] {
			continue
		}

		target := filepath.Join(parent, candidate)
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

func numberedNameParts(name string) (string, string, int) {
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	separator := strings.LastIndex(base, " ")
	if separator < 0 {
		return base, ext, 1
	}

	index, err := strconv.Atoi(base[separator+1:])
	if err != nil || index < 1 {
		return base, ext, 1
	}

	return base[:separator], ext, index + 1
}

func numberedName(base string, ext string, index int) string {
	return base + " " + strconv.Itoa(index) + ext
}
