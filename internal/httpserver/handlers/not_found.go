package handlers

import "github.com/gofiber/fiber/v2"

func (h Handlers) NotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error":  "not found",
		"method": c.Method(),
		"path":   c.Path(),
	})
}
