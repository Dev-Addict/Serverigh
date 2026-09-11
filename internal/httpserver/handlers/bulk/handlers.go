package bulk

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/filesystem"
)

type Handlers struct {
	Files                 *filesystem.Service
	WriteActionResult     WriteActionResultFunc
	WriteModeDisabled     WriteModeDisabledFunc
	WriteOperationalError WriteOperationalErrorFunc
	WriteRefresh          WriteRefreshFunc
}

type WriteActionResultFunc func(
	*fiber.Ctx,
	string,
	time.Time,
	error,
	...any,
) error

type WriteModeDisabledFunc func(*fiber.Ctx) (bool, error)

type WriteOperationalErrorFunc func(*fiber.Ctx, error) error

type WriteRefreshFunc func(*fiber.Ctx, string) error

func formPath(c *fiber.Ctx) string {
	return c.FormValue("path", "/")
}
