package handlers

import "github.com/gofiber/fiber/v2"

func (h Handlers) Root(c *fiber.Ctx) error {
	return c.Redirect("/browse?path=/", fiber.StatusFound)
}
