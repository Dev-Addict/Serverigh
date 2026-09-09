package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/config"
)

func TestBrowseRendersShell(t *testing.T) {
	root := t.TempDir()
	makeTestDir(t, root, "<script>")
	h := testHandlersWithRoot(t, root)
	app := fiber.New()
	app.Get("/browse", h.Browse)

	resp, body := testRequest(
		t,
		app,
		http.MethodGet,
		"/browse?path=%3Cscript%3E",
	)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, `class="app-shell"`) {
		t.Fatalf("expected browse shell in response, got %q", body)
	}

	if !strings.Contains(body, `hx-history-elt`) {
		t.Fatalf("expected htmx history element in response, got %q", body)
	}

	if !strings.Contains(body, `id="toast-region"`) {
		t.Fatalf("expected toast region in response, got %q", body)
	}

	if !strings.Contains(body, `id="search-region"`) {
		t.Fatalf("expected search region in response, got %q", body)
	}

	if !strings.Contains(body, `data-settings-open`) {
		t.Fatalf("expected settings button in response, got %q", body)
	}

	if !strings.Contains(body, `id="settings-modal"`) {
		t.Fatalf("expected settings modal in response, got %q", body)
	}

	if !strings.Contains(body, `id="keyboard-help-modal"`) {
		t.Fatalf("expected keyboard help modal in response, got %q", body)
	}

	if !strings.Contains(body, `id="write-modal"`) {
		t.Fatalf("expected write modal in response, got %q", body)
	}

	if !strings.Contains(body, `data-vim-key-status`) {
		t.Fatalf("expected vim key status in response, got %q", body)
	}

	if !strings.Contains(body, `data-help-open`) {
		t.Fatalf("expected keyboard help trigger in response, got %q", body)
	}

	if !strings.Contains(body, `Write Mode`) {
		t.Fatalf("expected write-mode shortcut section in response, got %q", body)
	}

	if !strings.Contains(body, `<option value="light" selected>Light</option>`) {
		t.Fatalf("expected light theme option in response, got %q", body)
	}

	if !strings.Contains(body, `<option value="dark">Dark</option>`) {
		t.Fatalf("expected dark theme option in response, got %q", body)
	}

	if strings.Contains(body, `id="search-region"><form`) {
		t.Fatalf("expected search control to be non-submittable, got %q", body)
	}

	if !strings.Contains(body, "No file selected") {
		t.Fatalf("expected empty preview state in response, got %q", body)
	}

	if strings.Contains(body, "<script>") {
		t.Fatalf("expected active path to be escaped, got %q", body)
	}

	if !strings.Contains(body, "&lt;script&gt;") {
		t.Fatalf("expected escaped active path in response, got %q", body)
	}
}

func TestBrowseShowsWriteMode(t *testing.T) {
	h := testHandlers(t, func(cfg *config.Config) {
		cfg.Write = true
	})
	app := fiber.New()
	app.Get("/browse", h.Browse)

	resp, body := testRequest(t, app, http.MethodGet, "/browse?path=/")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, "write-enabled") {
		t.Fatalf("expected write mode in response, got %q", body)
	}
}

func TestBrowseRestoresSelectedFilePreview(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "note.txt", "hello from preview")
	h := testHandlersWithRoot(t, root)
	app := fiber.New()
	app.Get("/browse", h.Browse)

	resp, body := testRequest(
		t,
		app,
		http.MethodGet,
		"/browse?path=/&file=/note.txt",
	)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, "hello from preview") {
		t.Fatalf("expected selected file preview in response, got %q", body)
	}

	if strings.Contains(body, "No file selected") {
		t.Fatalf("expected selected file to replace empty state, got %q", body)
	}
}

func TestBrowseIncludesHtmxNavigation(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "note.txt", "hello")
	h := testHandlersWithRoot(t, root)
	app := fiber.New()
	app.Get("/browse", h.Browse)

	resp, body := testRequest(t, app, http.MethodGet, "/browse?path=/")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, "/static/htmx.min.js") {
		t.Fatalf("expected htmx script in response, got %q", body)
	}

	if !strings.Contains(body, `hx-target="#preview-region"`) {
		t.Fatalf("expected preview htmx target in response, got %q", body)
	}

	if !strings.Contains(body, `href="/browse?`) ||
		!strings.Contains(body, `file=%2Fnote.txt`) ||
		!strings.Contains(body, `path=%2F`) {
		t.Fatalf("expected file link to preserve browse URL state, got %q", body)
	}

	if !strings.Contains(body, `hx-get="/preview?path=%2Fnote.txt"`) {
		t.Fatalf("expected file link to fetch preview partial, got %q", body)
	}
}

func TestBrowseUsesNormalizedPath(t *testing.T) {
	root := t.TempDir()
	makeTestDir(t, root, "docs")
	h := testHandlersWithRoot(t, root)
	app := fiber.New()
	app.Get("/browse", h.Browse)

	resp, body := testRequest(t, app, http.MethodGet, "/browse?path=/docs//")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if strings.Contains(body, "/docs//") {
		t.Fatalf("expected normalized path, got %q", body)
	}

	if !strings.Contains(body, "/docs") {
		t.Fatalf("expected normalized docs path, got %q", body)
	}
}
