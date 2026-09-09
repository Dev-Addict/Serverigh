package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h Handlers) Delete(c *fiber.Ctx) error {
	if disabled, err := h.writeModeDisabled(c); disabled {
		return err
	}

	start := time.Now()
	parentPath := formPath(c)
	target := c.FormValue("target")
	err := h.files.Trash(target)
	if resultErr := h.writeActionResult(
		c,
		"trash",
		start,
		err,
		"path",
		parentPath,
		"target",
		target,
	); resultErr != nil {
		return resultErr
	}

	return h.writeRefresh(c, parentPath)
}

func formPath(c *fiber.Ctx) string {
	return c.FormValue("path", "/")
}

func uploadName(filename string, relativePaths []string, index int) string {
	if index < len(relativePaths) && relativePaths[index] != "" {
		return relativePaths[index]
	}

	return filename
}
