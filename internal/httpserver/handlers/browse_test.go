package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/config"
)

func TestBrowseRendersShell(t *testing.T) {
	h := testHandlers(t)
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

	if !strings.Contains(body, "Serverigh is running.") {
		t.Fatalf("expected browse shell in response, got %q", body)
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
