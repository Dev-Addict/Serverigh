package write

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
)

func (h Handlers) Upload(c *fiber.Ctx) error {
	if disabled, err := h.ModeDisabled(c); disabled {
		return err
	}

	start := time.Now()
	form, err := c.MultipartForm()
	if err != nil {
		return h.ActionResult(c, "upload", start, filesystem.ErrInvalidPath)
	}

	parentPath := formPath(c)
	relativePaths := form.Value["relative_path"]
	fileHeaders := form.File["file"]
	if len(fileHeaders) == 0 {
		return h.ActionResult(
			c,
			"upload",
			start,
			filesystem.ErrInvalidPath,
			"path",
			parentPath,
		)
	}

	uploadPaths := make([]string, len(fileHeaders))
	for index, fileHeader := range fileHeaders {
		uploadPaths[index] = uploadName(fileHeader.Filename, relativePaths, index)
	}
	uploadPaths, err = h.Files.UniqueUploadPaths(parentPath, uploadPaths)
	if err != nil {
		return h.ActionResult(c, "upload", start, err, "path", parentPath)
	}

	for index, fileHeader := range fileHeaders {
		if err := h.uploadFile(parentPath, uploadPaths[index], fileHeader); err != nil {
			return h.ActionResult(c, "upload", start, err, "path", parentPath)
		}
	}
	h.LogAction(
		c,
		"upload",
		start,
		"path",
		parentPath,
		"file_count",
		len(fileHeaders),
	)

	return h.Refresh(c, parentPath)
}
