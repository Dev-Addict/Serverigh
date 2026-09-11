package bulk

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
)

func (h Handlers) bulkWrite(
	c *fiber.Ctx,
	action string,
	write bulkWriteFunc,
) error {
	if disabled, err := h.WriteModeDisabled(c); disabled {
		return err
	}

	targets, err := h.bulkResolvedTargets(c)
	if err != nil {
		return h.WriteOperationalError(c, err)
	}

	return h.bulkWriteTargets(c, action, targets, write)
}

func (h Handlers) bulkWriteTargets(
	c *fiber.Ctx,
	action string,
	targets []filesystem.ResolvedPath,
	write bulkWriteFunc,
) error {
	start := time.Now()
	parentPath := formPath(c)
	for _, target := range targets {
		if err := write(target); err != nil {
			return h.bulkWriteError(c, action, start, parentPath, target, targets, err)
		}
	}

	if err := h.WriteActionResult(
		c,
		action,
		start,
		nil,
		"path",
		parentPath,
		"target_count",
		len(targets),
	); err != nil {
		return err
	}

	return h.WriteRefresh(c, parentPath)
}

func (h Handlers) bulkWriteError(
	c *fiber.Ctx,
	action string,
	start time.Time,
	parentPath string,
	target filesystem.ResolvedPath,
	targets []filesystem.ResolvedPath,
	err error,
) error {
	return h.WriteActionResult(
		c,
		action,
		start,
		err,
		"path",
		parentPath,
		"target",
		target.Path,
		"target_count",
		len(targets),
	)
}
