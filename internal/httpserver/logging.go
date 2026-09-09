package httpserver

import (
	"log/slog"
	"runtime/debug"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
)

var requestCounter uint64

func registerMiddleware(app *fiber.App) {
	app.Use(requestLogMiddleware)
	app.Use(recoverMiddleware)
}

func recoverMiddleware(c *fiber.Ctx) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			slog.Error(
				"http panic recovered",
				"request_id", requestID(c),
				"method", c.Method(),
				"path", c.Path(),
				"panic", recovered,
				"stack", string(debug.Stack()),
			)
			err = c.Status(fiber.StatusInternalServerError).SendString(
				"internal server error",
			)
		}
	}()

	return c.Next()
}

func requestLogMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	id := nextRequestID()
	c.Locals("request_id", id)
	c.Set("X-Request-ID", id)

	err := c.Next()
	status := c.Response().StatusCode()
	if err != nil && status < fiber.StatusInternalServerError {
		status = fiber.StatusInternalServerError
	}
	level := slog.LevelInfo
	if status >= 500 {
		level = slog.LevelError
	} else if status >= 400 {
		level = slog.LevelWarn
	}

	slog.LogAttrs(
		c.UserContext(),
		level,
		"http request",
		slog.String("request_id", id),
		slog.String("method", c.Method()),
		slog.String("path", c.Path()),
		slog.Int("status", status),
		slog.Duration("duration", time.Since(start)),
		slog.String("ip", c.IP()),
		slog.String("user_agent", c.Get(fiber.HeaderUserAgent)),
		slog.Int("response_bytes", len(c.Response().Body())),
	)

	return err
}

func nextRequestID() string {
	id := atomic.AddUint64(&requestCounter, 1)

	return "req-" + strconv.FormatUint(id, 36)
}

func requestID(c *fiber.Ctx) string {
	if id, ok := c.Locals("request_id").(string); ok {
		return id
	}

	return ""
}
