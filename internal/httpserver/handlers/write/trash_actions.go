package write

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h Handlers) RestoreTrash(c *fiber.Ctx) error {
	if disabled, err := h.ModeDisabled(c); disabled {
		return err
	}

	start := time.Now()
	parentPath := formPath(c)
	target := c.FormValue("target")
	err := h.Files.RestoreTrash(target)
	if resultErr := h.ActionResult(
		c,
		"restore_trash",
		start,
		err,
		"path",
		parentPath,
		"target",
		target,
	); resultErr != nil {
		return resultErr
	}

	return h.Refresh(c, parentPath)
}

func (h Handlers) DeleteTrash(c *fiber.Ctx) error {
	if disabled, err := h.ModeDisabled(c); disabled {
		return err
	}

	start := time.Now()
	parentPath := formPath(c)
	target := c.FormValue("target")
	err := h.Files.DeleteTrash(target)
	if resultErr := h.ActionResult(
		c,
		"delete_trash",
		start,
		err,
		"path",
		parentPath,
		"target",
		target,
	); resultErr != nil {
		return resultErr
	}

	return h.Refresh(c, parentPath)
}
