package handlers

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/httpserver/view"
)

func (h Handlers) Preview(c *fiber.Ctx) error {
	start := time.Now()
	filePath := c.Query("path", "/")
	if c.Get("HX-Request") != "true" {
		return c.Redirect(
			view.BrowseURLWithFile(view.FolderForFile(filePath), filePath),
			fiber.StatusFound,
		)
	}

	preview, err := h.files.Preview(filePath)
	if err != nil {
		return h.writeOperationalError(c, err)
	}
	logHandlerOperation(
		c,
		slog.LevelInfo,
		"preview rendered",
		start,
		"path",
		preview.Path,
		"kind",
		preview.Kind,
		"mime_type",
		preview.MIMEType,
		"size",
		preview.Size,
		"bytes_read",
		preview.BytesRead,
		"truncated",
		preview.Truncated,
	)

	return h.render(c, "preview.html", preview)
}
