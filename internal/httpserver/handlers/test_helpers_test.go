package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/config"
)

func testHandlers(t *testing.T, opts ...func(*config.Config)) Handlers {
	t.Helper()

	cfg := config.Config{
		Root:            t.TempDir(),
		Host:            "127.0.0.1",
		Port:            4173,
		MaxPreviewBytes: 1024,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return New(cfg)
}

func testRequest(
	t *testing.T,
	app *fiber.App,
	method string,
	target string,
) (*http.Response, string) {
	t.Helper()

	req := httptest.NewRequest(method, target, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("test request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}

	return resp, string(body)
}
