package bulk

import (
	"io"
	"mime"

	"github.com/gofiber/fiber/v2"
)

func (h Handlers) BulkDownload(c *fiber.Ctx) error {
	targets, err := h.bulkResolvedTargets(c)
	if err != nil {
		return h.WriteOperationalError(c, err)
	}
	entries, err := collectBulkZipEntries(targets)
	if err != nil {
		return h.WriteOperationalError(c, err)
	}

	c.Set(fiber.HeaderXContentTypeOptions, "nosniff")
	c.Set(fiber.HeaderContentType, "application/zip")
	c.Set(
		fiber.HeaderContentDisposition,
		mime.FormatMediaType("attachment", map[string]string{
			"filename": "serverigh-selection.zip",
		}),
	)

	reader, writer := io.Pipe()
	go func() {
		writer.CloseWithError(h.writeBulkZip(writer, entries))
	}()

	return c.SendStream(reader)
}
