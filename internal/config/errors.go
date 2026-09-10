package config

import (
	"errors"
	"os"

	"serverigh/internal/apperror"
)

var (
	errInvalidTheme           = errors.New("theme must be light or dark")
	errInvalidMaxPreviewBytes = errors.New("max preview bytes must be positive")
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

func wrapConfigOperation(operation string, err error) error {
	return apperror.WrapOperation(
		apperror.CodeInvalidConfig,
		"invalid config",
		operation,
		err,
	)
}
