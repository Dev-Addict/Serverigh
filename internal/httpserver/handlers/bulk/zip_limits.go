package bulk

import "serverigh/internal/apperror"

const (
	maxBulkArchiveFiles = 2_000
	maxBulkArchiveBytes = int64(1 << 30)
)

func validateBulkZipLimits(entries []bulkZipEntry) error {
	if len(entries) > maxBulkArchiveFiles {
		return apperror.New(
			apperror.CodeInvalidPath,
			"too many files selected",
		)
	}

	var size int64
	for _, entry := range entries {
		size += entry.Size
		if size > maxBulkArchiveBytes {
			return apperror.New(
				apperror.CodeInvalidPath,
				"selected files are too large",
			)
		}
	}

	return nil
}
