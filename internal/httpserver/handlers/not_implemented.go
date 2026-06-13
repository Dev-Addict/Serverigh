package handlers

import "github.com/gofiber/fiber/v2"

func (h Handlers) NotImplemented(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"error":  "not implemented",
		"method": c.Method(),
		"path":   c.Path(),
	})
}
