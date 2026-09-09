package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestFoldersRendersInitialFolderBranch(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "docs", "drafts"), 0o755); err != nil {
		t.Fatalf("create folders: %v", err)
	}
	h := writeModeHandlers(t, root)
	app := fiber.New()
	app.Get("/partials/folders", h.Folders)

	req := httptest.NewRequest(
		http.MethodGet,
		"/partials/folders?selected=/docs/drafts",
		nil,
	)
	resp, body := testFiberRequest(t, app, req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.StatusCode, body)
	}
	for _, expected := range []string{
		`data-folder-path="/"`,
		`data-folder-path="/docs"`,
		`data-folder-path="/docs/drafts"`,
		`data-selected="true"`,
	} {
		if !strings.Contains(body, expected) {
			t.Fatalf("expected body to contain %q: %s", expected, body)
		}
	}
}

func TestFoldersRendersOneChildLevel(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "docs", "drafts", "deep"), 0o755); err != nil {
		t.Fatalf("create folders: %v", err)
	}
	h := writeModeHandlers(t, root)
	app := fiber.New()
	app.Get("/partials/folders", h.Folders)

	req := httptest.NewRequest(http.MethodGet, "/partials/folders?path=/docs", nil)
	resp, body := testFiberRequest(t, app, req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.StatusCode, body)
	}
	if !strings.Contains(body, `data-folder-path="/docs/drafts"`) {
		t.Fatalf("expected child folder in response, got %s", body)
	}
	if strings.Contains(body, `data-folder-path="/docs/drafts/deep"`) {
		t.Fatalf("expected one child level only, got %s", body)
	}
}
