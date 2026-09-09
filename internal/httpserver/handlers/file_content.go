package handlers

import (
	"log/slog"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
)

func (h Handlers) Download(c *fiber.Ctx) error {
	start := time.Now()
	file, err := h.files.Open(c.Query("path", "/"))
	if err != nil {
		return h.writeOperationalError(c, err)
	}
	logFileStream(c, "download started", start, file)

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

func (h Handlers) Raw(c *fiber.Ctx) error {
	start := time.Now()
	file, err := h.files.Open(c.Query("path", "/"))
	if err != nil {
		return h.writeOperationalError(c, err)
	}
	logFileStream(c, "raw file started", start, file)

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

func sendOpenedFile(c *fiber.Ctx, file filesystem.OpenedFile) error {
	size := streamSize(file.Size)
	if size < 0 {
		if err := c.SendStream(file.Handle); err != nil {
			file.Close()

			return err
		}

		return nil
	}

	if err := c.SendStream(file.Handle, size); err != nil {
		file.Close()

		return err
	}

	return nil
}

func streamSize(size int64) int {
	maxInt := int64(int(^uint(0) >> 1))
	if size < 0 || size > maxInt {
		return -1
	}

	return int(size)
}

func logFileStream(
	c *fiber.Ctx,
	event string,
	start time.Time,
	file filesystem.OpenedFile,
) {
	logHandlerOperation(
		c,
		slog.LevelInfo,
		event,
		start,
		"path",
		file.Path,
		"mime_type",
		file.MIMEType,
		"size",
		file.Size,
	)
}
