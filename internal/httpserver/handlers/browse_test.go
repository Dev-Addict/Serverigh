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
