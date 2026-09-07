package searcher

import (
	"errors"
	"io/fs"

	"serverigh/internal/apperror"
)

var (
	ErrInvalidPath = apperror.New(
		apperror.CodeInvalidPath,
		"invalid path",
	)
	ErrOutsideRoot = apperror.New(
		apperror.CodeOutsideRoot,
		"path escapes configured root",
	)
)

func wrapFileError(operation string, err error) error {
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return apperror.WrapOperation(
			apperror.CodeNotFound,
			"path not found",
			operation,
			err,
		)
	case errors.Is(err, fs.ErrPermission):
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
