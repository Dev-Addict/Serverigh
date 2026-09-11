package write

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h Handlers) Move(c *fiber.Ctx) error {
	if disabled, err := h.ModeDisabled(c); disabled {
		return err
	}

	start := time.Now()
	parentPath := formPath(c)
	target := c.FormValue("target")
	destination := c.FormValue("destination", "/")
	err := h.Files.Move(target, destination)
	if resultErr := h.logTransferResult(
		c,
		"move",
		start,
		err,
		parentPath,
		target,
		destination,
	); resultErr != nil {
		return resultErr
	}

	return h.Refresh(c, parentPath)
}

func (h Handlers) Copy(c *fiber.Ctx) error {
	if disabled, err := h.ModeDisabled(c); disabled {
		return err
	}

	start := time.Now()
	parentPath := formPath(c)
	target := c.FormValue("target")
	destination := c.FormValue("destination", parentPath)
	err := h.Files.Copy(target, destination)
	if resultErr := h.logTransferResult(
		c,
		"copy",
		start,
		err,
		parentPath,
		target,
		destination,
	); resultErr != nil {
		return resultErr
	}

	return h.Refresh(c, parentPath)
}

func (h Handlers) Duplicate(c *fiber.Ctx) error {
	if disabled, err := h.ModeDisabled(c); disabled {
		return err
	}

	start := time.Now()
	parentPath := formPath(c)
	target := c.FormValue("target")
	err := h.Files.Duplicate(target)
	if resultErr := h.ActionResult(
		c,
		"duplicate",
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

func (h Handlers) logTransferResult(
	c *fiber.Ctx,
	action string,
	start time.Time,
	err error,
	parentPath string,
	target string,
	destination string,
) error {
	return h.ActionResult(
		c,
		action,
		start,
		err,
		"path",
		parentPath,
		"target",
		target,
		"destination",
		destination,
	)
}
