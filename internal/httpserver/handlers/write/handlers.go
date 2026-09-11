package write

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/apperror"
	"serverigh/internal/filesystem"
)

type Handlers struct {
	Config                Config
	Files                 *filesystem.Service
	Render                RenderFunc
	WriteError            WriteErrorFunc
	WriteOperationalError WriteOperationalErrorFunc
}

type Config struct {
	Root    string
	Enabled bool
}

type RenderFunc func(*fiber.Ctx, string, any) error

type WriteErrorFunc func(
	*fiber.Ctx,
	int,
	apperror.Code,
	string,
) error

type WriteOperationalErrorFunc func(*fiber.Ctx, error) error

type ActionResultFunc func(
	*fiber.Ctx,
	string,
	time.Time,
	error,
	...any,
) error

type ModeDisabledFunc func(*fiber.Ctx) (bool, error)

type RefreshFunc func(*fiber.Ctx, string) error

func formPath(c *fiber.Ctx) string {
	return c.FormValue("path", "/")
}

func uploadName(filename string, relativePaths []string, index int) string {
	if index < len(relativePaths) && relativePaths[index] != "" {
		return relativePaths[index]
	}

	return filename
}
