package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestSettingsPersistsLocalConfig(t *testing.T) {
	root := t.TempDir()
	h := testHandlersWithRoot(t, root)
	app := fiber.New()
	app.Post("/settings", h.Settings)

	form := url.Values{}
	form.Set("theme", "dark")
	form.Set("max_preview_bytes", "2048")
	form.Set("column_size", "true")
	form.Set("column_modified", "false")
	form.Set("column_created", "true")
	form.Set("column_mode", "false")

	resp, _ := testFormRequest(t, app, "/settings", form)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	content, err := os.ReadFile(
		filepath.Join(root, ".serverigh", "config.toml"),
	)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	for _, expected := range []string{
		`theme = "dark"`,
		`max_preview_bytes = 2048`,
		`modified = false`,
		`mode = false`,
	} {
		if !strings.Contains(string(content), expected) {
			t.Fatalf("expected %q in config:\n%s", expected, content)
		}
	}

	if h.config.Theme != "dark" {
		t.Fatalf("expected in-memory theme to update, got %q", h.config.Theme)
	}
	if h.config.MaxPreviewBytes != 2048 {
		t.Fatalf(
			"expected in-memory max preview bytes 2048, got %d",
			h.config.MaxPreviewBytes,
		)
	}
	if h.config.Columns.Modified {
		t.Fatalf("expected in-memory modified column to be disabled")
	}
	if h.config.Columns.Mode {
		t.Fatalf("expected in-memory mode column to be disabled")
	}
}

func testFormRequest(
	t *testing.T,
	app *fiber.App,
	target string,
	form url.Values,
) (*http.Response, string) {
	t.Helper()

	reqBody := strings.NewReader(form.Encode())
	req := httptest.NewRequest(http.MethodPost, target, reqBody)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
