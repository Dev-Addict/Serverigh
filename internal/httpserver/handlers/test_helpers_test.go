package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

	handlers, err := New(cfg)
	if err != nil {
		t.Fatalf("create handlers: %v", err)
	}

	return handlers
}

func testHandlersWithRoot(
	t *testing.T,
	root string,
	opts ...func(*config.Config),
) Handlers {
	t.Helper()

	cfg := config.Config{
		Root:            root,
		Host:            "127.0.0.1",
		Port:            4173,
		MaxPreviewBytes: 1024,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	handlers, err := New(cfg)
	if err != nil {
		t.Fatalf("create handlers: %v", err)
	}

	return handlers
}

func testRequest(
	t *testing.T,
	app *fiber.App,
	method string,
	target string,
) (*http.Response, string) {
	t.Helper()

	return testRequestWithHeaders(t, app, method, target, nil)
}

func testRequestWithHeaders(
	t *testing.T,
	app *fiber.App,
	method string,
	target string,
	headers map[string]string,
) (*http.Response, string) {
	t.Helper()

	req := httptest.NewRequest(method, target, nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

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

func writeTestFile(t *testing.T, root string, name string, content string) {
	t.Helper()

	filename := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		t.Fatalf("create test directory: %v", err)
	}

	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}
}

func makeTestDir(t *testing.T, root string, name string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Join(root, name), 0o755); err != nil {
		t.Fatalf("create test directory: %v", err)
	}
}
