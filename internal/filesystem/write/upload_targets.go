package write

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func (s Service) uploadTarget(parentPath string, uploadPath string) (string, error) {
	cleanPath, err := cleanUploadPath(uploadPath)
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

	target, err := s.ensureUploadDirectories(parent.Absolute, path.Dir(cleanPath))
	if err != nil {
		return "", err
	}

	return s.childTargetForAbsolute(target, path.Base(cleanPath))
}

func cleanUploadPath(uploadPath string) (string, error) {
	slashedPath := strings.ReplaceAll(uploadPath, "\\", "/")
	if slashedPath == "" || strings.ContainsRune(slashedPath, '\x00') {
		return "", ErrInvalidPath
	}

	segments := strings.Split(strings.TrimPrefix(slashedPath, "/"), "/")
	for _, segment := range segments {
		if err := validateEntryName(segment); err != nil {
			return "", err
		}
	}

	cleanPath := path.Clean("/" + strings.TrimPrefix(slashedPath, "/"))
	if cleanPath == "/" {
		return "", ErrInvalidPath
	}

	return strings.TrimPrefix(cleanPath, "/"), nil
}

func (s Service) ensureUploadDirectories(root string, directory string) (string, error) {
	current := root
	if directory == "." {
		return current, nil
	}

	for _, segment := range strings.Split(directory, "/") {
		next := filepath.Join(current, segment)
		if err := s.ensureInsideRoot(next); err != nil {
			return "", err
		}
		if err := ensureUploadDirectory(next); err != nil {
			return "", err
		}
		current = next
	}

	return current, nil
}

func ensureUploadDirectory(target string) error {
	info, err := os.Lstat(target)
	if errors.Is(err, os.ErrNotExist) {
		return os.Mkdir(target, 0o755)
	}
	if err != nil {
		return wrapFileError("inspect upload directory", err)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return ErrNotDirectory
	}

	return nil
}
