package bulk

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"serverigh/internal/apperror"
	"serverigh/internal/filesystem"
)

func validateBulkTransfer(
	targets []filesystem.ResolvedPath,
	destination filesystem.ResolvedPath,
	replace bool,
) error {
	names := map[string]bool{}
	for _, target := range targets {
		if err := validateBulkTransferTarget(target, destination, replace); err != nil {
			return err
		}
		name := path.Base(target.Path)
		if !replace && names[name] {
			return apperror.New(
				apperror.CodeAlreadyExists,
				"path already exists",
			)
		}
		names[name] = true
	}

	return nil
}

func validateBulkTransferTarget(
	target filesystem.ResolvedPath,
	destination filesystem.ResolvedPath,
	replace bool,
) error {
	if target.Info.IsDir() && sameOrDescendant(destination.Absolute, target.Absolute) {
		return apperror.New(
			apperror.CodeInvalidPath,
			"cannot place directory inside itself",
		)
	}
	if replace {
		return nil
	}

	return validateMoveTargetAvailable(target, destination)
}

func validateMoveTargetAvailable(
	target filesystem.ResolvedPath,
	destination filesystem.ResolvedPath,
) error {
	filename := filepath.Join(destination.Absolute, path.Base(target.Path))
	if _, err := os.Lstat(filename); err == nil {
		return apperror.New(apperror.CodeAlreadyExists, "path already exists")
	} else if err != nil && !os.IsNotExist(err) {
		return apperror.WrapOperation(
			apperror.CodeFilesystem,
			"filesystem error",
			"inspect target path",
			err,
		)
	}

	return nil
}

func sameOrDescendant(candidate string, parent string) bool {
	relative, err := filepath.Rel(parent, candidate)
	if err != nil {
		return false
	}

	return relative == "." ||
		(relative != ".." && !strings.HasPrefix(
			relative,
			".."+string(filepath.Separator),
		))
}
