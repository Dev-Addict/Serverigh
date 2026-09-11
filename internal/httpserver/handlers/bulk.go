package handlers

import (
	"github.com/gofiber/fiber/v2"

	bulkhandlers "serverigh/internal/httpserver/handlers/bulk"
)

func (h Handlers) BulkDelete(c *fiber.Ctx) error {
	return h.bulkHandlers().BulkDelete(c)
}

func (h Handlers) BulkDuplicate(c *fiber.Ctx) error {
	return h.bulkHandlers().BulkDuplicate(c)
}

func (h Handlers) BulkCopy(c *fiber.Ctx) error {
	return h.bulkHandlers().BulkCopy(c)
}

func (h Handlers) BulkMove(c *fiber.Ctx) error {
	return h.bulkHandlers().BulkMove(c)
}

func (h Handlers) BulkDownload(c *fiber.Ctx) error {
	return h.bulkHandlers().BulkDownload(c)
}

func (h Handlers) bulkHandlers() bulkhandlers.Handlers {
	writeHandlers := h.writeHandlers()

	return bulkhandlers.Handlers{
		Files:                 h.files,
		WriteActionResult:     writeHandlers.ActionResult,
		WriteModeDisabled:     writeHandlers.ModeDisabled,
		WriteOperationalError: h.writeOperationalError,
		WriteRefresh:          writeHandlers.Refresh,
	}
}
