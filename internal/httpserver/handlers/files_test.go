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
		t.Fatalf("expected status 200, got %d: %s", resp.StatusCode, body)
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
		"/browse?",
	) || !strings.Contains(
		strings.ToLower(body),
		"file=%2fa%3fb%23c%26d.txt",
	) || !strings.Contains(
		strings.ToLower(body),
		"path=%2f",
	) {
		t.Fatalf("expected escaped browse file path, got %q", body)
	}

	if !strings.Contains(
		strings.ToLower(body),
		`hx-get="/preview?path=%2fa%3fb%23c%26d.txt"`,
	) {
		t.Fatalf("expected escaped htmx preview path, got %q", body)
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
		`id="search-region"`,
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

func TestFilesAppliesTableControlQuery(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "alpha.txt", "a")
	writeTestFile(t, root, "beta.txt", "bbb")
	writeTestFile(t, root, ".env", "hidden")
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/partials/files", h.Files)

	resp, body := testRequestWithHeaders(
		t,
		app,
		http.MethodGet,
		"/partials/files?path=/&sort=size&dir=desc&filter=bet&hidden=1",
		map[string]string{"HX-Request": "true"},
	)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, "beta.txt") {
		t.Fatalf("expected file in response, got %q", body)
	}

	if !strings.Contains(body, "alpha.txt") {
		t.Fatalf("expected second file in response, got %q", body)
	}

	if strings.Contains(body, ".env") {
		t.Fatalf("expected hidden file to follow config, got %q", body)
	}

	location := resp.Header.Get("HX-Push-Url")
	for _, part := range []string{
		"/browse?",
		"sort=size",
		"dir=desc",
	} {
		if !strings.Contains(location, part) {
			t.Fatalf("expected %q in pushed URL %q", part, location)
		}
	}

	for _, removed := range []string{"filter=", "hidden="} {
		if strings.Contains(location, removed) {
			t.Fatalf("expected %q to be removed from pushed URL %q", removed, location)
		}
	}
}

func TestFilesRendersTableControls(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "note.txt", "hello")
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/partials/files", h.Files)

	resp, body := testRequest(
		t,
		app,
		http.MethodGet,
		"/partials/files?path=/&sort=modified&dir=desc&hidden=1",
	)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	expected := []string{
		`aria-sort="descending"`,
		`data-entry-row`,
		`data-relative-path="note.txt"`,
		`data-absolute-path=`,
		`<time datetime=`,
		`data-copy-text="note.txt"`,
		`title="Copy absolute path"`,
	}

	for _, part := range expected {
		if !strings.Contains(body, part) {
			t.Fatalf("expected %q in response, got %q", part, body)
		}
	}

	for _, removed := range []string{
		`>Kind<`,
		`sort=kind`,
		`name="filter"`,
		`filter=`,
		`class="table-controls"`,
		`data-hidden-toggle`,
		`name="hidden"`,
		`>Hidden files<`,
	} {
		if strings.Contains(body, removed) {
			t.Fatalf("expected %q to be removed from table, got %q", removed, body)
		}
	}
}
