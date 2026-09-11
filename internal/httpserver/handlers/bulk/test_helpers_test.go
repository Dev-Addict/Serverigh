package bulk

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/apperror"
	"serverigh/internal/filesystem"
)

func testHandlers(t *testing.T) Handlers {
	t.Helper()

	return testHandlersWithRoot(t, t.TempDir())
}

func testHandlersWithRoot(t *testing.T, root string) Handlers {
	t.Helper()

	return testHandlersForRoot(t, root, false)
}

func writeModeHandlers(t *testing.T, root string) Handlers {
	t.Helper()

	return testHandlersForRoot(t, root, true)
}

func testHandlersForRoot(t *testing.T, root string, write bool) Handlers {
	t.Helper()

	files, err := filesystem.New(root, true, 1024)
	if err != nil {
		t.Fatalf("create filesystem: %v", err)
	}

	return Handlers{
		Files:                 &files,
		WriteActionResult:     testWriteActionResult,
		WriteModeDisabled:     testWriteModeDisabled(write),
		WriteOperationalError: testWriteOperationalError,
		WriteRefresh: func(c *fiber.Ctx, _ string) error {
			return c.SendStatus(fiber.StatusOK)
		},
	}
}

func testWriteActionResult(
	c *fiber.Ctx,
	_ string,
	_ time.Time,
	err error,
	_ ...any,
) error {
	if err != nil {
		return testWriteOperationalError(c, err)
	}

	return nil
}

func testWriteModeDisabled(write bool) WriteModeDisabledFunc {
	return func(c *fiber.Ctx) (bool, error) {
		if write {
			return false, nil
		}

		return true, c.SendStatus(fiber.StatusForbidden)
	}
}

func testWriteOperationalError(c *fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	if opErr, ok := apperror.AsOperational(err); ok {
		status = testStatusForOperationalCode(opErr.Code)
	}

	return c.SendStatus(status)
}

func testStatusForOperationalCode(code apperror.Code) int {
	switch code {
	case apperror.CodeInvalidPath:
		return fiber.StatusBadRequest
	case apperror.CodeNotDirectory:
		return fiber.StatusBadRequest
	case apperror.CodeNotFound:
		return fiber.StatusNotFound
	case apperror.CodeAlreadyExists:
		return fiber.StatusConflict
	case apperror.CodePermissionDenied:
		return fiber.StatusForbidden
	default:
		return fiber.StatusInternalServerError
	}
}

func formRequest(
	t *testing.T,
	app *fiber.App,
	target string,
	values url.Values,
) (*http.Response, string) {
	t.Helper()

	req := httptest.NewRequest(
		http.MethodPost,
		target,
		strings.NewReader(values.Encode()),
	)
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationForm)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("test request: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}
	if err := resp.Body.Close(); err != nil {
		t.Fatalf("close response body: %v", err)
	}

	return resp, string(body)
}

func writeTestFile(
	t *testing.T,
	root string,
	relativePath string,
	content string,
) {
	t.Helper()

	filename := filepath.Join(root, relativePath)
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		t.Fatalf("create parent directory: %v", err)
	}
	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}
}
