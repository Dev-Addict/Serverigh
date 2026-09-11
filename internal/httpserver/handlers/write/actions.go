package write

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h Handlers) Delete(c *fiber.Ctx) error {
	if disabled, err := h.ModeDisabled(c); disabled {
		return err
	}

	start := time.Now()
	parentPath := formPath(c)
	target := c.FormValue("target")
	err := h.Files.Trash(target)
	if resultErr := h.ActionResult(
		c,
		"trash",
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
