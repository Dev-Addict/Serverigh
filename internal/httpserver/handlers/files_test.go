package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestFilesRendersDirectoryListing(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "note.txt", "hello")
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/partials/files", h.Files)

	resp, body := testRequest(
		t,
		app,
		http.MethodGet,
		"/partials/files?path=/",
	)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, "note.txt") {
		t.Fatalf("expected listed file in response, got %q", body)
	}
}

func TestFilesRejectsTraversal(t *testing.T) {
	h := testHandlers(t)
	app := fiber.New()
	app.Get("/partials/files", h.Files)

	resp, _ := testRequest(
		t,
		app,
		http.MethodGet,
		"/partials/files?path=../outside",
	)

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", resp.StatusCode)
	}
}

func TestFilesReturnsOperationalErrorCode(t *testing.T) {
	h := testHandlers(t)
	app := fiber.New()
	app.Get("/partials/files", h.Files)

	resp, body := testRequest(
		t,
		app,
		http.MethodGet,
		"/partials/files?path=../outside",
	)

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, `"code":"outside_root"`) {
		t.Fatalf("expected outside root code, got %q", body)
	}
}

func TestFilesEscapesEntryLinks(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "a?b#c&d.txt", "hello")
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/partials/files", h.Files)

	resp, body := testRequest(
		t,
		app,
		http.MethodGet,
		"/partials/files?path=/",
	)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(
		strings.ToLower(body),
		"/preview?path=%2fa%3fb%23c%26d.txt",
	) {
		t.Fatalf("expected escaped preview path, got %q", body)
	}
}

func TestFilesReturnsHtmlErrorForHtmx(t *testing.T) {
	h := testHandlers(t)
	app := fiber.New()
	app.Get("/partials/files", h.Files)

	resp, body := testRequestWithHeaders(
		t,
		app,
		http.MethodGet,
		"/partials/files?path=../outside",
		map[string]string{"HX-Request": "true"},
	)

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get(fiber.HeaderContentType)
	if !strings.Contains(contentType, "text/html") {
		t.Fatalf("expected html response, got %q", contentType)
	}

	if !strings.Contains(body, `data-error-code="outside_root"`) {
		t.Fatalf("expected rendered operational error, got %q", body)
	}
}

func TestFilesHtmxResponseIncludesOutOfBandRegions(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "note.txt", "hello")
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/partials/files", h.Files)

	resp, body := testRequestWithHeaders(
		t,
		app,
		http.MethodGet,
		"/partials/files?path=/",
		map[string]string{"HX-Request": "true"},
	)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	expected := []string{
		`id="path-summary"`,
		`id="breadcrumbs-region"`,
		`id="preview-region"`,
		`id="status-row"`,
		`hx-swap-oob=`,
	}

	for _, part := range expected {
		if !strings.Contains(body, part) {
			t.Fatalf("expected %q in htmx response, got %q", part, body)
		}
	}
}
