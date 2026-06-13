package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestHealthReturnsStatus(t *testing.T) {
	h := testHandlers(t)
	app := fiber.New()
	app.Get("/healthz", h.Health)

	resp, body := testRequest(t, app, http.MethodGet, "/healthz")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatalf("decode health response: %v", err)
	}

	if payload["status"] != "ok" {
		t.Fatalf("expected health status ok, got %#v", payload["status"])
	}

	if payload["app"] != "serverigh" {
		t.Fatalf("expected app serverigh, got %#v", payload["app"])
	}
}
