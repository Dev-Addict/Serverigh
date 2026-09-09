package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
)

func (h Handlers) Upload(c *fiber.Ctx) error {
	if disabled, err := h.writeModeDisabled(c); disabled {
		return err
	}

	start := time.Now()
	form, err := c.MultipartForm()
	if err != nil {
		return h.writeActionResult(c, "upload", start, filesystem.ErrInvalidPath)
	}

	parentPath := formPath(c)
	relativePaths := form.Value["relative_path"]
	fileHeaders := form.File["file"]
	if len(fileHeaders) == 0 {
		return h.writeActionResult(
			c,
			"upload",
			start,
			filesystem.ErrInvalidPath,
			"path",
			parentPath,
		)
	}

	for index, fileHeader := range fileHeaders {
		if err := h.uploadFile(parentPath, relativePaths, index, fileHeader); err != nil {
			return h.writeActionResult(c, "upload", start, err, "path", parentPath)
		}
	}
	h.logWriteAction(
		c,
		"upload",
		start,
		"path",
		parentPath,
		"file_count",
		len(fileHeaders),
	)

	return h.writeRefresh(c, parentPath)
}
