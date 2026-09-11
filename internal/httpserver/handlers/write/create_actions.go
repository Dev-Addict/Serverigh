package write

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h Handlers) Mkdir(c *fiber.Ctx) error {
	if disabled, err := h.ModeDisabled(c); disabled {
		return err
	}

	start := time.Now()
	parentPath := formPath(c)
	err := h.Files.CreateDirectory(parentPath, c.FormValue("name"))
	if resultErr := h.ActionResult(
		c,
		"mkdir",
		start,
		err,
		"path",
		parentPath,
	); resultErr != nil {
		return resultErr
	}

	return h.Refresh(c, parentPath)
}

func (h Handlers) CreateFile(c *fiber.Ctx) error {
	if disabled, err := h.ModeDisabled(c); disabled {
		return err
	}

	start := time.Now()
	parentPath := formPath(c)
	err := h.Files.CreateFile(parentPath, c.FormValue("name"))
	if resultErr := h.ActionResult(
		c,
		"create_file",
		start,
		err,
		"path",
		parentPath,
	); resultErr != nil {
		return resultErr
	}

	return h.Refresh(c, parentPath)
}

func (h Handlers) Rename(c *fiber.Ctx) error {
	if disabled, err := h.ModeDisabled(c); disabled {
		return err
	}

	start := time.Now()
	parentPath := formPath(c)
	requestPath := c.FormValue("target")
	err := h.Files.Rename(requestPath, c.FormValue("name"))
	if resultErr := h.ActionResult(
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

	return h.Refresh(c, parentPath)
}
