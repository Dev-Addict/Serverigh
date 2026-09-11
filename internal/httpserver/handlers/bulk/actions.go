package bulk

import (
	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
)

type bulkWriteFunc func(filesystem.ResolvedPath) error

func (h Handlers) BulkDelete(c *fiber.Ctx) error {
	return h.bulkWrite(c, "bulk_trash", func(target filesystem.ResolvedPath) error {
		return h.Files.Trash(target.Path)
	})
}

func (h Handlers) BulkDuplicate(c *fiber.Ctx) error {
	return h.bulkWrite(
		c,
		"bulk_duplicate",
		func(target filesystem.ResolvedPath) error {
			return h.Files.Duplicate(target.Path)
		},
	)
}

func (h Handlers) BulkCopy(c *fiber.Ctx) error {
	if disabled, err := h.WriteModeDisabled(c); disabled {
		return err
	}

	destination, err := h.bulkDestination(c.FormValue("destination", formPath(c)))
	if err != nil {
		return h.WriteOperationalError(c, err)
	}

	return h.bulkTransfer(c, "bulk_copy", destination, true, func(target filesystem.ResolvedPath) error {
		return h.Files.Copy(target.Path, destination.Path)
	})
}

func (h Handlers) BulkMove(c *fiber.Ctx) error {
	if disabled, err := h.WriteModeDisabled(c); disabled {
		return err
	}

	destination, err := h.bulkDestination(c.FormValue("destination", "/"))
	if err != nil {
		return h.WriteOperationalError(c, err)
	}

	return h.bulkTransfer(c, "bulk_move", destination, false, func(target filesystem.ResolvedPath) error {
		return h.Files.Move(target.Path, destination.Path)
	})
}
