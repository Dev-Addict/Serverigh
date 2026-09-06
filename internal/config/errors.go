package config

import (
	"errors"
	"os"

	"serverigh/internal/apperror"
)

func wrapRootOperation(operation string, err error) error {
	switch {
	case errors.Is(err, os.ErrNotExist):
		return apperror.WrapOperation(
			apperror.CodeNotFound,
			"root path not found",
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
