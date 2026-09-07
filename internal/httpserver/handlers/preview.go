package handlers

import (
	"github.com/gofiber/fiber/v2"

	"serverigh/internal/httpserver/view"
)

func (h Handlers) Preview(c *fiber.Ctx) error {
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

	return h.render(c, "preview.html", preview)
}
