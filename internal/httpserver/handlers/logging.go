package handlers

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func logHandlerOperation(
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
