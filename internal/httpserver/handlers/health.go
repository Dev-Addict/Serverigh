package handlers

import "github.com/gofiber/fiber/v2"

func (h Handlers) Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"app":                "serverigh",
		"status":             "ok",
		"root":               h.config.Root,
		"write":              h.config.Write,
		"showHidden":         h.config.ShowHidden,
		"maxPreviewBytes":    h.config.MaxPreviewBytes,
		"filesystemBrowsing": "pending",
	})
}
