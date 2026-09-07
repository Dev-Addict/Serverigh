package filesystem

import (
	"os"
	"path/filepath"

	"serverigh/internal/apperror"
	"serverigh/internal/filesystem/searcher"
)

type Service struct {
	root            string
	showHidden      bool
	maxPreviewBytes int64
	searchIndex     *searcher.Index
}

func New(root string, showHidden bool, maxPreviewBytes int64) (Service, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return Service{}, apperror.WrapOperation(
			apperror.CodeInvalidPath,
			"invalid root path",
			"resolve root path",
			err,
		)
	}

	realRoot, err := filepath.EvalSymlinks(absRoot)
	if err != nil {
		return Service{}, wrapFileError("resolve root symlinks", err)
	}

	info, err := os.Stat(realRoot)
	if err != nil {
		return Service{}, wrapFileError("inspect root path", err)
	}

	if !info.IsDir() {
		return Service{}, ErrNotDirectory
	}

	if maxPreviewBytes < 1 {
		return Service{}, apperror.New(
			apperror.CodeInvalidConfig,
			"max preview bytes must be greater than zero",
		)
	}

	return Service{
		root:            realRoot,
		showHidden:      showHidden,
		maxPreviewBytes: maxPreviewBytes,
		searchIndex:     searcher.NewIndex(),
	}, nil
}

func (s Service) Root() string {
	return s.root
}
