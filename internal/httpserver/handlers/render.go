package handlers

import (
	"bytes"

	"github.com/gofiber/fiber/v2"
)

func (h Handlers) render(c *fiber.Ctx, name string, data any) error {
	var body bytes.Buffer
	if err := h.templates.ExecuteTemplate(&body, name, data); err != nil {
		return err
	}

	return c.Type("html", "utf-8").Send(body.Bytes())
}
