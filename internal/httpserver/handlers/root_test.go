package handlers

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRootRedirectsToBrowse(t *testing.T) {
	h := testHandlers(t)
	app := fiber.New()
	app.Get("/", h.Root)

	resp, _ := testRequest(t, app, http.MethodGet, "/")

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("expected status 302, got %d", resp.StatusCode)
	}

	if location := resp.Header.Get("Location"); location != "/browse?path=/" {
		t.Fatalf("expected browse redirect, got %q", location)
	}
}
