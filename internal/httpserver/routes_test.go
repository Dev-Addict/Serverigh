package httpserver

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

func TestRoutesWireCoreHandlers(t *testing.T) {
	cfg := testConfig(t)
	if err := os.Mkdir(filepath.Join(cfg.Root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs directory: %v", err)
	}

	if err := os.WriteFile(
		filepath.Join(cfg.Root, "note.txt"),
		[]byte("hello"),
		0o644,
	); err != nil {
		t.Fatalf("create note file: %v", err)
	}

	app := testApp(t, cfg)

	tests := []struct {
		name   string
		path   string
		status int
	}{
		{
			name:   "root",
			path:   "/",
			status: http.StatusFound,
		},
		{
			name:   "browse",
			path:   "/browse?path=/docs",
			status: http.StatusOK,
		},
		{
			name:   "health",
			path:   "/healthz",
			status: http.StatusOK,
		},
		{
			name:   "files",
			path:   "/partials/files?path=/",
			status: http.StatusOK,
		},
		{
			name:   "preview",
			path:   "/preview?path=/note.txt",
			status: http.StatusOK,
		},
		{
			name:   "raw",
			path:   "/raw?path=/note.txt",
			status: http.StatusOK,
		},
		{
			name:   "download",
			path:   "/download?path=/note.txt",
			status: http.StatusOK,
		},
		{
			name:   "static",
			path:   "/static/app.js",
			status: http.StatusOK,
		},
		{
			name:   "unknown",
			path:   "/missing",
			status: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("test request: %v", err)
			}

			if resp.StatusCode != tt.status {
				t.Fatalf("expected status %d, got %d", tt.status, resp.StatusCode)
			}
		})
	}
}

func TestRoutesServeEmbeddedStaticAssets(t *testing.T) {
	app := testApp(t, testConfig(t))
	req := httptest.NewRequest(http.MethodGet, "/static/app.js", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("test request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}

	if string(body) == "" {
		t.Fatalf("expected static asset body")
	}
}

func TestNewReturnsSetupErrors(t *testing.T) {
	cfg := testConfig(t)
	cfg.Root = filepath.Join(t.TempDir(), "missing")

	if _, err := New(cfg); err == nil {
		t.Fatalf("expected setup error")
	}
}

func testApp(t *testing.T, cfg config.Config) *fiber.App {
	t.Helper()

	app, err := New(cfg)
	if err != nil {
		t.Fatalf("create app: %v", err)
	}

	return app
}

func testConfig(t *testing.T) config.Config {
	t.Helper()

	return config.Config{
		Root:            t.TempDir(),
		Host:            "127.0.0.1",
		Port:            4173,
		MaxPreviewBytes: 1024,
	}
}
