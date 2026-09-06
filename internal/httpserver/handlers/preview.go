package handlers

import "github.com/gofiber/fiber/v2"

func (h Handlers) Preview(c *fiber.Ctx) error {
	preview, err := h.files.Preview(c.Query("path", "/"))
	if err != nil {
		return h.writeOperationalError(c, err)
	}

	return h.render(c, "preview.html", preview)
}
