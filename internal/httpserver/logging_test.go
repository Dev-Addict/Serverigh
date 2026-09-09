package httpserver

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRequestLogMiddlewareAddsRequestID(t *testing.T) {
	var output bytes.Buffer
	original := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&output, nil)))
	t.Cleanup(func() {
		slog.SetDefault(original)
	})

	app := fiber.New()
	registerMiddleware(app)
	app.Get("/ok", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/ok?token=secret", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("test request: %v", err)
	}

	if resp.Header.Get("X-Request-ID") == "" {
		t.Fatalf("expected request id header")
	}

	logs := output.String()
	if !strings.Contains(logs, "msg=\"http request\"") {
		t.Fatalf("expected request log, got %q", logs)
	}
	if strings.Contains(logs, "token=secret") {
		t.Fatalf("expected request log to omit query string, got %q", logs)
	}
}

func TestRecoverMiddlewareLogsPanic(t *testing.T) {
	var output bytes.Buffer
	original := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&output, nil)))
	t.Cleanup(func() {
		slog.SetDefault(original)
	})

	app := fiber.New()
	registerMiddleware(app)
	app.Get("/panic", func(c *fiber.Ctx) error {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("test request: %v", err)
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", resp.StatusCode)
	}
	if !strings.Contains(output.String(), "msg=\"http panic recovered\"") {
		t.Fatalf("expected panic log, got %q", output.String())
	}
}

func TestNextRequestIDIsUnique(t *testing.T) {
	first := nextRequestID()
	second := nextRequestID()

	if first == second {
		t.Fatalf("expected unique request ids, got %q", first)
	}
}
