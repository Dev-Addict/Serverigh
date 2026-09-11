package bulk

import (
	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
)

func (h Handlers) bulkTransfer(
	c *fiber.Ctx,
	action string,
	destination filesystem.ResolvedPath,
	replace bool,
	write bulkWriteFunc,
) error {
	targets, err := h.bulkResolvedTargets(c)
	if err != nil {
		return h.WriteOperationalError(c, err)
	}
	if err := validateBulkTransfer(targets, destination, replace); err != nil {
		return h.WriteOperationalError(c, err)
	}

	return h.bulkWriteTargets(c, action, targets, write)
}
