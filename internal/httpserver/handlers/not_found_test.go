package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestNotFoundReturnsRouteDetails(t *testing.T) {
	h := testHandlers(t)
	app := fiber.New()
	app.Use(h.NotFound)

	resp, body := testRequest(t, app, http.MethodGet, "/missing")

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, `"error":"not found"`) {
		t.Fatalf("expected not found error in response, got %q", body)
	}

	if !strings.Contains(body, `"path":"/missing"`) {
		t.Fatalf("expected route path in response, got %q", body)
	}
}
