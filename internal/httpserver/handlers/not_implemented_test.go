package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestNotImplementedReturnsRouteDetails(t *testing.T) {
	h := testHandlers(t)
	app := fiber.New()
	app.Get("/pending", h.NotImplemented)

	resp, body := testRequest(t, app, http.MethodGet, "/pending")

	if resp.StatusCode != http.StatusNotImplemented {
		t.Fatalf("expected status 501, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, `"error":"not implemented"`) {
		t.Fatalf("expected not implemented error in response, got %q", body)
	}

	if !strings.Contains(body, `"path":"/pending"`) {
		t.Fatalf("expected route path in response, got %q", body)
	}
}
