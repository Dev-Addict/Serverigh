package handlers

import (
	"mime"

	"github.com/gofiber/fiber/v2"
)

func (h Handlers) Download(c *fiber.Ctx) error {
	file, err := h.files.Open(c.Query("path", "/"))
	if err != nil {
		return h.writeOperationalError(c, err)
	}

	c.Set(fiber.HeaderXContentTypeOptions, "nosniff")
	c.Set(
		fiber.HeaderContentDisposition,
		mime.FormatMediaType("attachment", map[string]string{
			"filename": file.Name,
		}),
	)

	if file.MIMEType != "" {
		c.Set(fiber.HeaderContentType, file.MIMEType)
	}

	return sendOpenedFile(c, file)
}
