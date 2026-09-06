package filesystem

import (
	"errors"
	"os"

	"serverigh/internal/apperror"
)

var (
	ErrInvalidPath = apperror.New(
		apperror.CodeInvalidPath,
		"invalid path",
	)
	ErrOutsideRoot = apperror.New(
		apperror.CodeOutsideRoot,
		"path is outside configured root",
	)
	ErrIsDirectory = apperror.New(
		apperror.CodeIsDirectory,
		"path is a directory",
	)
	ErrNotDirectory = apperror.New(
		apperror.CodeNotDirectory,
		"path is not a directory",
	)
	ErrPathChanged = apperror.New(
		apperror.CodeInvalidPath,
		"path changed while opening",
	)
)

func wrapFileError(operation string, err error) error {
	switch {
	case errors.Is(err, os.ErrNotExist):
		return apperror.WrapOperation(
			apperror.CodeNotFound,
			"path not found",
			operation,
			err,
		)
	case errors.Is(err, os.ErrPermission):
		return apperror.WrapOperation(
			apperror.CodePermissionDenied,
			"permission denied",
			operation,
			err,
		)
	default:
		return apperror.WrapOperation(
			apperror.CodeFilesystem,
			"filesystem error",
			operation,
			err,
		)
	}
}
