package handlers

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/apperror"
)

type ErrorView struct {
	Code    string
	Message string
}

func (h Handlers) writeOperationalError(c *fiber.Ctx, err error) error {
	if errors.Is(err, context.Canceled) {
		return h.writeError(
			c,
			fiber.StatusRequestTimeout,
			apperror.CodeRequestCanceled,
			"request canceled",
		)
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return h.writeError(
			c,
			fiber.StatusGatewayTimeout,
			apperror.CodeRequestTimeout,
			"request timed out",
		)
	}

	opErr, ok := apperror.AsOperational(err)
	if !ok {
		return h.writeError(
			c,
			fiber.StatusInternalServerError,
			apperror.CodeFilesystem,
			"filesystem error",
		)
	}

	return h.writeError(
		c,
		statusForOperationalCode(opErr.Code),
		opErr.Code,
		opErr.Message,
	)
}

func (h Handlers) writeError(
	c *fiber.Ctx,
	status int,
	code apperror.Code,
	message string,
) error {
	logHandlerError(c, status, code, message)
	if wantsHTMLError(c) {
		c.Status(status)

		return h.render(c, "error.html", ErrorView{
			Code:    string(code),
			Message: message,
		})
	}

	return c.Status(status).JSON(fiber.Map{
		"code":  string(code),
		"error": message,
	})
}

func logHandlerError(
	c *fiber.Ctx,
	status int,
	code apperror.Code,
	message string,
) {
	level := slog.LevelWarn
	if status >= fiber.StatusInternalServerError {
		level = slog.LevelError
	}

	slog.Log(
		c.UserContext(),
		level,
		"request error",
		"request_id",
		c.Locals("request_id"),
		"method",
		c.Method(),
		"path",
		c.Path(),
		"status",
		status,
		"code",
		code,
		"message",
		message,
	)
}

func wantsHTMLError(c *fiber.Ctx) bool {
	if c.Get("HX-Request") == "true" {
		return true
	}

	return strings.Contains(c.Get(fiber.HeaderAccept), "text/html")
}

func statusForOperationalCode(code apperror.Code) int {
	status := fiber.StatusInternalServerError

	switch code {
	case apperror.CodeInvalidPath:
		status = fiber.StatusBadRequest
	case apperror.CodeInvalidConfig:
		status = fiber.StatusBadRequest
	case apperror.CodeOutsideRoot:
		status = fiber.StatusForbidden
	case apperror.CodeIsDirectory:
		status = fiber.StatusBadRequest
	case apperror.CodeNotDirectory:
		status = fiber.StatusBadRequest
	case apperror.CodeNotFound:
		status = fiber.StatusNotFound
	case apperror.CodeAlreadyExists:
		status = fiber.StatusConflict
	case apperror.CodePermissionDenied:
		status = fiber.StatusForbidden
	case apperror.CodeRequestCanceled:
		status = fiber.StatusRequestTimeout
	case apperror.CodeRequestTimeout:
		status = fiber.StatusGatewayTimeout
	case apperror.CodeWriteDisabled:
		status = fiber.StatusForbidden
	case apperror.CodeServer:
		status = fiber.StatusInternalServerError
	}

	return status
}
