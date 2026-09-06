package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/config"
)

func TestPreviewRendersBoundedFilePreview(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "note.txt", "abcdef")
	h := testHandlersWithRoot(t, root, func(cfg *config.Config) {
		cfg.MaxPreviewBytes = 3
	})

	app := fiber.New()
	app.Get("/preview", h.Preview)

	resp, body := testRequest(t, app, http.MethodGet, "/preview?path=/note.txt")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, "abc") {
		t.Fatalf("expected preview content, got %q", body)
	}

	if strings.Contains(body, "abcdef") {
		t.Fatalf("expected bounded preview, got %q", body)
	}

	if !strings.Contains(body, "Preview truncated") {
		t.Fatalf("expected truncation notice, got %q", body)
	}
}

func TestPreviewRejectsDirectory(t *testing.T) {
	h := testHandlers(t)
	app := fiber.New()
	app.Get("/preview", h.Preview)

	resp, _ := testRequest(t, app, http.MethodGet, "/preview?path=/")

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestPreviewReturnsOperationalErrorCode(t *testing.T) {
	h := testHandlers(t)
	app := fiber.New()
	app.Get("/preview", h.Preview)

	resp, body := testRequest(t, app, http.MethodGet, "/preview?path=/")

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, `"code":"is_directory"`) {
		t.Fatalf("expected is directory code, got %q", body)
	}
}

func TestPreviewEscapesActionLinks(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "a?b#c&d.txt", "hello")
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/preview", h.Preview)

	resp, body := testRequest(
		t,
		app,
		http.MethodGet,
		"/preview?path=%2Fa%3Fb%23c%26d.txt",
	)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	normalizedBody := strings.ToLower(body)
	if !strings.Contains(normalizedBody, "/raw?path=%2fa%3fb%23c%26d.txt") {
		t.Fatalf("expected escaped raw path, got %q", body)
	}

	if !strings.Contains(
		normalizedBody,
		"/download?path=%2fa%3fb%23c%26d.txt",
	) {
		t.Fatalf("expected escaped download path, got %q", body)
	}
}
