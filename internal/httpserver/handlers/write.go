package handlers

import (
	"github.com/gofiber/fiber/v2"

	writehandlers "serverigh/internal/httpserver/handlers/write"
)

func (h Handlers) Mkdir(c *fiber.Ctx) error {
	return h.writeHandlers().Mkdir(c)
}

func (h Handlers) CreateFile(c *fiber.Ctx) error {
	return h.writeHandlers().CreateFile(c)
}

func (h Handlers) Rename(c *fiber.Ctx) error {
	return h.writeHandlers().Rename(c)
}

func (h Handlers) Move(c *fiber.Ctx) error {
	return h.writeHandlers().Move(c)
}

func (h Handlers) Copy(c *fiber.Ctx) error {
	return h.writeHandlers().Copy(c)
}

func (h Handlers) Duplicate(c *fiber.Ctx) error {
	return h.writeHandlers().Duplicate(c)
}

func (h Handlers) Delete(c *fiber.Ctx) error {
	return h.writeHandlers().Delete(c)
}

func (h Handlers) RestoreTrash(c *fiber.Ctx) error {
	return h.writeHandlers().RestoreTrash(c)
}

func (h Handlers) DeleteTrash(c *fiber.Ctx) error {
	return h.writeHandlers().DeleteTrash(c)
}

func (h Handlers) Upload(c *fiber.Ctx) error {
	return h.writeHandlers().Upload(c)
}

func (h Handlers) Folders(c *fiber.Ctx) error {
	return h.writeHandlers().Folders(c)
}

func (h Handlers) Write() writehandlers.Handlers {
	return h.writeHandlers()
}

func (h Handlers) writeHandlers() writehandlers.Handlers {
	return writehandlers.Handlers{
		Config: writehandlers.Config{
			Root:    h.config.Root,
			Enabled: h.config.Write,
		},
		Files:                 h.files,
		Render:                h.render,
		WriteError:            h.writeError,
		WriteOperationalError: h.writeOperationalError,
	}
}
