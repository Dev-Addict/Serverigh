package handlers

import "github.com/gofiber/fiber/v2"

func (h Handlers) Files(c *fiber.Ctx) error {
	listing, err := h.files.List(c.Query("path", "/"))
	if err != nil {
		return h.writeOperationalError(c, err)
	}

	return h.render(c, "files.html", listing)
}
