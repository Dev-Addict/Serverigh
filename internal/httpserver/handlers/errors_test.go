package handlers

import (
	"testing"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/apperror"
)

func TestStatusForOperationalCode(t *testing.T) {
	tests := []struct {
		code   apperror.Code
		status int
	}{
		{
			code:   apperror.CodeInvalidPath,
			status: fiber.StatusBadRequest,
		},
		{
			code:   apperror.CodeInvalidConfig,
			status: fiber.StatusBadRequest,
		},
		{
			code:   apperror.CodeOutsideRoot,
			status: fiber.StatusForbidden,
		},
		{
			code:   apperror.CodeNotFound,
			status: fiber.StatusNotFound,
		},
		{
			code:   apperror.CodePermissionDenied,
			status: fiber.StatusForbidden,
		},
		{
			code:   apperror.CodeFilesystem,
			status: fiber.StatusInternalServerError,
		},
		{
			code:   apperror.CodeRequestCanceled,
			status: fiber.StatusRequestTimeout,
		},
		{
			code:   apperror.CodeRequestTimeout,
			status: fiber.StatusGatewayTimeout,
		},
		{
			code:   apperror.CodeServer,
			status: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.code), func(t *testing.T) {
			got := statusForOperationalCode(tt.code)
			if got != tt.status {
				t.Fatalf("expected status %d, got %d", tt.status, got)
			}
		})
	}
}
