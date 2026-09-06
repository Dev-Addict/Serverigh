package handlers

import (
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
)

func (h Handlers) Raw(c *fiber.Ctx) error {
	file, err := h.files.Open(c.Query("path", "/"))
	if err != nil {
		return h.writeOperationalError(c, err)
	}

	c.Set(fiber.HeaderXContentTypeOptions, "nosniff")
	c.Set(fiber.HeaderContentSecurityPolicy, "default-src 'none'")
	c.Set(fiber.HeaderContentType, rawContentType(file.File))

	return sendOpenedFile(c, file)
}

func rawContentType(file filesystem.File) string {
	if isActiveRawContent(file) {
		return "text/plain; charset=utf-8"
	}

	if file.MIMEType == "" {
		return "application/octet-stream"
	}

	return file.MIMEType
}

func isActiveRawContent(file filesystem.File) bool {
	switch strings.ToLower(filepath.Ext(file.Name)) {
	case ".htm", ".html", ".js", ".mjs", ".svg", ".xhtml", ".xml":
		return true
	}

	mimeType := strings.ToLower(strings.TrimSpace(file.MIMEType))
	mimeType, _, _ = strings.Cut(mimeType, ";")

	switch mimeType {
	case "application/javascript",
		"application/xhtml+xml",
		"application/xml",
		"image/svg+xml",
		"text/html",
		"text/javascript",
		"text/xml":
		return true
	default:
		return false
	}
}
