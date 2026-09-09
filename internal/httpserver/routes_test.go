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
	cfg.Write = true
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
		header map[string]string
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
			name:   "folders",
			path:   "/partials/folders?path=/",
			status: http.StatusOK,
		},
		{
			name:   "breadcrumbs",
			path:   "/partials/breadcrumbs?path=/docs",
			status: http.StatusOK,
		},
		{
			name:   "search",
			path:   "/partials/search?path=/&q=note",
			status: http.StatusOK,
		},
		{
			name:   "preview",
			path:   "/preview?path=/note.txt",
			header: map[string]string{"HX-Request": "true"},
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
			for key, value := range tt.header {
				req.Header.Set(key, value)
			}

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

	tests := []string{
		"/static/app.js",
		"/static/app/help.js",
		"/static/app/keyboard.js",
		"/static/app/keyboard-command-state.js",
		"/static/app/keyboard-copy-path.js",
		"/static/app/keyboard-dom.js",
		"/static/app/keyboard-entries.js",
		"/static/app/keyboard-events.js",
		"/static/app/keyboard-htmx-navigation.js",
		"/static/app/keyboard-navigation.js",
		"/static/app/keyboard-panes.js",
		"/static/app/keyboard-runner.js",
		"/static/app/keyboard-shortcuts.js",
		"/static/app/keyboard-table.js",
		"/static/app/keyboard-write-actions.js",
		"/static/app/layout.js",
		"/static/app/layout-preview.js",
		"/static/app/theme.js",
		"/static/app/vim-status.js",
		"/static/app/write-actions.js",
		"/static/app/write-context-menu.js",
		"/static/app/write-context-state.js",
		"/static/app/write-create.js",
		"/static/app/write-entry-actions.js",
		"/static/app/write-folder-cache.js",
		"/static/app/write-folder-dom.js",
		"/static/app/write-folder-requests.js",
		"/static/app/write-menu.js",
		"/static/app/write-modal-elements.js",
		"/static/app/write-modal-folders.js",
		"/static/app/write-modal-view.js",
		"/static/app/write-modal.js",
		"/static/app/write-request.js",
		"/static/app/write-upload.js",
	}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
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
		})
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
