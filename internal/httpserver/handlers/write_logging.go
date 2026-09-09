package handlers

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h Handlers) writeActionResult(
	c *fiber.Ctx,
	action string,
	start time.Time,
	err error,
	attrs ...any,
) error {
	if err != nil {
		h.logWriteActionAt(c, slog.LevelWarn, action, start, attrs...)

		return h.writeOperationalError(c, err)
	}

	h.logWriteAction(c, action, start, attrs...)

	return nil
}

func (h Handlers) logWriteAction(
	c *fiber.Ctx,
	action string,
	start time.Time,
	attrs ...any,
) {
	h.logWriteActionAt(c, slog.LevelInfo, action, start, attrs...)
}

func (h Handlers) logWriteActionAt(
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
