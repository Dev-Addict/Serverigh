package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/apperror"
)

type ErrorView struct {
	Code    string
	Message string
}

func (h Handlers) writeOperationalError(c *fiber.Ctx, err error) error {
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
	case apperror.CodePermissionDenied:
		status = fiber.StatusForbidden
	case apperror.CodeServer:
		status = fiber.StatusInternalServerError
	}

	return status
}
