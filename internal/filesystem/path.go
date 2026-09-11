package filesystem

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"serverigh/internal/apperror"
)

type ResolvedPath struct {
	Path     string
	Absolute string
	Info     os.FileInfo
}

func (s Service) Resolve(requestPath string) (ResolvedPath, error) {
	logicalPath, err := cleanRequestPath(requestPath)
	if err != nil {
		return ResolvedPath{}, err
	}

	if isTrashPath(logicalPath) {
		return s.resolveTrashPath(logicalPath)
	}

	absolutePath := filepath.Join(
		s.root,
		filepath.FromSlash(strings.TrimPrefix(logicalPath, "/")),
	)

	if err := s.ensureInsideRoot(absolutePath); err != nil {
		return ResolvedPath{}, err
	}

	realPath, err := filepath.EvalSymlinks(absolutePath)
	if err != nil {
		return ResolvedPath{}, wrapFileError("resolve path symlinks", err)
	}

	if err := s.ensureInsideRoot(realPath); err != nil {
		return ResolvedPath{}, err
	}

	info, err := os.Stat(realPath)
	if err != nil {
		return ResolvedPath{}, wrapFileError("inspect path", err)
	}

	return ResolvedPath{
		Path:     logicalPath,
		Absolute: realPath,
		Info:     info,
	}, nil
}

func (s Service) resolveTrashPath(logicalPath string) (ResolvedPath, error) {
	absolutePath := trashStorageChild(s.root, logicalPath)
	if err := s.ensureInsideRoot(absolutePath); err != nil {
		return ResolvedPath{}, err
	}

	realPath, err := filepath.EvalSymlinks(absolutePath)
	if err != nil {
		return ResolvedPath{}, wrapFileError("resolve path symlinks", err)
	}

	if err := s.ensureInsideRoot(realPath); err != nil {
		return ResolvedPath{}, err
	}

	info, err := os.Stat(realPath)
	if err != nil {
		return ResolvedPath{}, wrapFileError("inspect path", err)
	}

	return ResolvedPath{
		Path:     logicalPath,
		Absolute: realPath,
		Info:     info,
	}, nil
}

func cleanRequestPath(requestPath string) (string, error) {
	if strings.ContainsRune(requestPath, '\x00') {
		return "", ErrInvalidPath
	}

	if requestPath == "" {
		requestPath = "/"
	}

	slashedPath := strings.ReplaceAll(requestPath, "\\", "/")
	for _, segment := range strings.Split(slashedPath, "/") {
		if segment == ".." {
			return "", ErrOutsideRoot
		}
	}

	cleanPath := path.Clean("/" + strings.TrimPrefix(slashedPath, "/"))
	if cleanPath == "." {
		return "/", nil
	}

	return cleanPath, nil
}

func (s Service) ensureInsideRoot(candidate string) error {
	rel, err := filepath.Rel(s.root, candidate)
	if err != nil {
		return apperror.WrapOperation(
			apperror.CodeFilesystem,
			"filesystem error",
			"compare path with root",
			err,
		)
	}

	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return ErrOutsideRoot
	}

	return nil
}

func displayPath(parentPath string, name string) string {
	if parentPath == "/" {
		return "/" + name
	}

	return path.Join(parentPath, name)
}
