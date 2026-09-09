package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h Handlers) Move(c *fiber.Ctx) error {
	if disabled, err := h.writeModeDisabled(c); disabled {
		return err
	}

	start := time.Now()
	parentPath := formPath(c)
	target := c.FormValue("target")
	destination := c.FormValue("destination", "/")
	err := h.files.Move(target, destination)
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

	return h.writeRefresh(c, parentPath)
}

func (h Handlers) Copy(c *fiber.Ctx) error {
	if disabled, err := h.writeModeDisabled(c); disabled {
		return err
	}

	start := time.Now()
	parentPath := formPath(c)
	target := c.FormValue("target")
	destination := c.FormValue("destination", parentPath)
	err := h.files.Copy(target, destination)
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

	return h.writeRefresh(c, parentPath)
}

func (h Handlers) Duplicate(c *fiber.Ctx) error {
	if disabled, err := h.writeModeDisabled(c); disabled {
		return err
	}

	start := time.Now()
	parentPath := formPath(c)
	target := c.FormValue("target")
	err := h.files.Duplicate(target)
	if resultErr := h.writeActionResult(
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

	return h.writeRefresh(c, parentPath)
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
	return h.writeActionResult(
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
