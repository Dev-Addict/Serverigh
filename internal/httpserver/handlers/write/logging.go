package write

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h Handlers) ActionResult(
	c *fiber.Ctx,
	action string,
	start time.Time,
	err error,
	attrs ...any,
) error {
	if err != nil {
		h.logActionAt(c, slog.LevelWarn, action, start, attrs...)

		return h.WriteOperationalError(c, err)
	}

	h.LogAction(c, action, start, attrs...)

	return nil
}

func (h Handlers) LogAction(
	c *fiber.Ctx,
	action string,
	start time.Time,
	attrs ...any,
) {
	h.logActionAt(c, slog.LevelInfo, action, start, attrs...)
}

func (h Handlers) logActionAt(
	c *fiber.Ctx,
	level slog.Level,
	action string,
	start time.Time,
	attrs ...any,
) {
	values := []any{
		"request_id", c.Locals("request_id"),
		"action", action,
		"duration", time.Since(start),
	}
	values = append(values, attrs...)

	slog.Log(c.UserContext(), level, "write action", values...)
}

func logOperation(
	c *fiber.Ctx,
	level slog.Level,
	event string,
	start time.Time,
	attrs ...any,
) {
	values := []any{
		"request_id", c.Locals("request_id"),
		"duration", time.Since(start),
	}
	values = append(values, attrs...)

	slog.Log(c.UserContext(), level, event, values...)
}
