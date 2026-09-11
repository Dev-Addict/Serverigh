package write

import (
	"os"
	"path"
)

func (s Service) Move(requestPath string, targetParentPath string) error {
	source, err := s.resolve(requestPath)
	if err != nil {
		return err
	}

	parent, err := s.resolve(targetParentPath)
	if err != nil {
		return err
	}
	if !parent.Info.IsDir() {
		return ErrNotDirectory
	}
	if source.Info.IsDir() && sameOrDescendant(parent.Absolute, source.Absolute) {
		return ErrInvalidPath
	}

	target, err := s.childTargetForAbsolute(parent.Absolute, path.Base(source.Path))
	if err != nil {
		return err
	}

	if err := os.Rename(source.Absolute, target); err != nil {
		return wrapWriteError("move path", err)
	}
	if isTrashPath(source.Path) {
		return s.removeTrashMetadata(path.Base(source.Path))
	}

	return nil
}

func (s Service) Copy(requestPath string, targetParentPath string) error {
	source, err := s.resolve(requestPath)
	if err != nil {
		return err
	}

	parent, err := s.resolve(targetParentPath)
	if err != nil {
		return err
	}
	if !parent.Info.IsDir() {
		return ErrNotDirectory
	}
	if source.Info.IsDir() && sameOrDescendant(parent.Absolute, source.Absolute) {
		return ErrInvalidPath
	}

	target, err := s.uniqueChildTarget(parent.Absolute, path.Base(source.Path))
	if err != nil {
		return err
	}

	if source.Info.IsDir() {
		err = copyDirectory(source.Absolute, target)
	} else {
		err = copyFile(source.Absolute, target, source.Info.Mode())
	}
	if err != nil {
		return wrapWriteError("copy path", err)
	}

	return nil
}
