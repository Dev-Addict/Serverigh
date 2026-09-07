package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestSearchRendersRankedResults(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "docs/report.txt", "docs")
	writeTestFile(t, root, "archive/report.txt", "archive")
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/partials/search", h.Search)

	resp, body := testRequest(
		t,
		app,
		http.MethodGet,
		"/partials/search?path=/docs&q=report",
	)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, `class="file-panel search-panel"`) {
		t.Fatalf("expected search results in files panel, got %q", body)
	}

	if !strings.Contains(body, `class="file-table search-results-table"`) {
		t.Fatalf("expected search results table, got %q", body)
	}

	if strings.Contains(body, `aria-sort=`) {
		t.Fatalf("expected search table headers to be unsortable, got %q", body)
	}

	docsIndex := strings.Index(body, "path=%2Fdocs")
	archiveIndex := strings.Index(body, "path=%2Farchive")
	if docsIndex == -1 || archiveIndex == -1 {
		t.Fatalf("expected both search results, got %q", body)
	}

	if docsIndex > archiveIndex {
		t.Fatalf("expected current directory result first, got %q", body)
	}
}

func TestSearchEmptyQueryRestoresDirectoryListing(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "note.txt", "hello")
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/partials/search", h.Search)

	resp, body := testRequest(
		t,
		app,
		http.MethodGet,
		"/partials/search?path=/&q=&sort=modified&dir=desc",
	)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	expected := []string{
		`class="file-panel"`,
		"note.txt",
		`aria-sort="descending"`,
	}
	for _, part := range expected {
		if !strings.Contains(body, part) {
			t.Fatalf("expected %q in directory listing, got %q", part, body)
		}
	}
}

func TestSearchFileResultLinksToParentAndPreview(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "docs/report.txt", "docs")
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/partials/search", h.Search)

	resp, body := testRequest(
		t,
		app,
		http.MethodGet,
		"/partials/search?path=/&q=report",
	)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	expected := `href="/browse?path=%2Fdocs&amp;file=%2Fdocs%2Freport.txt"`
	if !strings.Contains(body, expected) {
		t.Fatalf("expected file result to open parent and preview, got %q", body)
	}

	for _, part := range []string{
		`data-search-preview-link`,
		`data-preview-url="/preview?path=%2Fdocs%2Freport.txt"`,
		`<th scope="col">Actions</th>`,
		`class="path-actions"`,
		`title="Go">Go</a>`,
	} {
		if !strings.Contains(body, part) {
			t.Fatalf("expected %q in file result, got %q", part, body)
		}
	}

	for _, removed := range []string{
		`data-copy-text=`,
	} {
		if strings.Contains(body, removed) {
			t.Fatalf("expected %q to be removed from file result, got %q", removed, body)
		}
	}
}

func TestSearchDirectoryResultLinksToDirectory(t *testing.T) {
	root := t.TempDir()
	makeTestDir(t, root, "docs")
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/partials/search", h.Search)

	resp, body := testRequest(
		t,
		app,
		http.MethodGet,
		"/partials/search?path=/&q=docs",
	)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, `href="/browse?path=%2Fdocs"`) {
		t.Fatalf("expected directory result to open directory, got %q", body)
	}
}

func TestSearchRejectsInvalidActivePath(t *testing.T) {
	h := testHandlers(t)

	app := fiber.New()
	app.Get("/partials/search", h.Search)

	resp, body := testRequest(
		t,
		app,
		http.MethodGet,
		"/partials/search?path=../outside&q=note",
	)

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, `"code":"outside_root"`) {
		t.Fatalf("expected outside root code, got %q", body)
	}
}
