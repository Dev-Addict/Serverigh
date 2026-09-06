package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRawReturnsWholeFile(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "note.txt", "hello")
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/raw", h.Raw)

	resp, body := testRequest(t, app, http.MethodGet, "/raw?path=/note.txt")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if body != "hello" {
		t.Fatalf("expected whole file body, got %q", body)
	}
}

func TestRawDowngradesActiveContent(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "index.html", "<script>alert(1)</script>")
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/raw", h.Raw)

	resp, body := testRequest(t, app, http.MethodGet, "/raw?path=/index.html")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get(fiber.HeaderContentType)
	if !strings.Contains(contentType, "text/plain") {
		t.Fatalf("expected safe text content type, got %q", contentType)
	}

	if resp.Header.Get(fiber.HeaderXContentTypeOptions) != "nosniff" {
		t.Fatalf("expected nosniff header")
	}

	if resp.Header.Get(fiber.HeaderContentSecurityPolicy) != "default-src 'none'" {
		t.Fatalf("expected restrictive content security policy")
	}

	if body != "<script>alert(1)</script>" {
		t.Fatalf("expected raw body, got %q", body)
	}
}
