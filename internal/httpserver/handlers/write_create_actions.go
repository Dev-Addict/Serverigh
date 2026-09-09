package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h Handlers) Mkdir(c *fiber.Ctx) error {
	if disabled, err := h.writeModeDisabled(c); disabled {
		return err
	}

	start := time.Now()
	parentPath := formPath(c)
	err := h.files.CreateDirectory(parentPath, c.FormValue("name"))
	if resultErr := h.writeActionResult(
		c,
		"mkdir",
		start,
		err,
		"path",
		parentPath,
	); resultErr != nil {
		return resultErr
	}

	return h.writeRefresh(c, parentPath)
}

func (h Handlers) CreateFile(c *fiber.Ctx) error {
	if disabled, err := h.writeModeDisabled(c); disabled {
		return err
	}

	start := time.Now()
	parentPath := formPath(c)
	err := h.files.CreateFile(parentPath, c.FormValue("name"))
	if resultErr := h.writeActionResult(
		c,
		"create_file",
		start,
		err,
		"path",
		parentPath,
	); resultErr != nil {
		return resultErr
	}

	return h.writeRefresh(c, parentPath)
}

func (h Handlers) Rename(c *fiber.Ctx) error {
	if disabled, err := h.writeModeDisabled(c); disabled {
		return err
	}

	start := time.Now()
	parentPath := formPath(c)
	requestPath := c.FormValue("target")
	err := h.files.Rename(requestPath, c.FormValue("name"))
	if resultErr := h.writeActionResult(
		c,
		"rename",
		start,
		err,
		"path",
		parentPath,
		"target",
		requestPath,
	); resultErr != nil {
		return resultErr
	}

	return h.writeRefresh(c, parentPath)
}
