package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestDownloadReturnsAttachment(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "note.txt", "hello")
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/download", h.Download)

	resp, body := testRequest(t, app, http.MethodGet, "/download?path=/note.txt")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	disposition := resp.Header.Get("Content-Disposition")
	if !strings.Contains(disposition, "attachment") {
		t.Fatalf("expected attachment disposition, got %q", disposition)
	}

	if !strings.Contains(disposition, "note.txt") {
		t.Fatalf("expected filename in disposition, got %q", disposition)
	}

	if resp.Header.Get(fiber.HeaderXContentTypeOptions) != "nosniff" {
		t.Fatalf("expected nosniff header")
	}

	if body != "hello" {
		t.Fatalf("expected downloaded body, got %q", body)
	}
}
