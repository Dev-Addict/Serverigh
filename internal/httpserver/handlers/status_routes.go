package handlers

import "github.com/gofiber/fiber/v2"

func (h Handlers) Root(c *fiber.Ctx) error {
	return c.Redirect("/browse?path=/", fiber.StatusFound)
}

func (h Handlers) Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"app":                "serverigh",
		"status":             "ok",
		"root":               h.config.Root,
		"theme":              h.config.Theme,
		"write":              h.config.Write,
		"showHidden":         h.config.ShowHidden,
		"maxPreviewBytes":    h.config.MaxPreviewBytes,
		"filesystemBrowsing": "ready",
	})
}

func (h Handlers) NotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error":  "not found",
		"method": c.Method(),
		"path":   c.Path(),
	})
}

func (h Handlers) NotImplemented(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"error":  "not implemented",
		"method": c.Method(),
		"path":   c.Path(),
	})
}
