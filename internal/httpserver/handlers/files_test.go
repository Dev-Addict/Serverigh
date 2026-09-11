package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
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

func decodeFilesVersion(t *testing.T, body string) filesVersionResponse {
	t.Helper()
	var response filesVersionResponse
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		t.Fatalf("decode files version response %q: %v", body, err)
	}

	return response
}

func filesVersionFor(
	t *testing.T,
	app *fiber.App,
	path string,
) filesVersionResponse {
	t.Helper()
	resp, body := testRequest(
		t,
		app,
		http.MethodGet,
		"/partials/files/version?path="+path,
	)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.StatusCode, body)
	}

	return decodeFilesVersion(t, body)
}

func TestFilesVersionChangesWhenDirectoryContentsChange(t *testing.T) {
	root := t.TempDir()
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/partials/files/version", h.FilesVersion)

	resp, body := testRequest(
		t,
		app,
		http.MethodGet,
		"/partials/files/version?path=/",
	)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.StatusCode, body)
	}
	before := decodeFilesVersion(t, body)

	writeTestFile(t, root, "note.txt", "hello")
	resp, body = testRequest(
		t,
		app,
		http.MethodGet,
		"/partials/files/version?path=/",
	)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.StatusCode, body)
	}
	after := decodeFilesVersion(t, body)

	if after.Path != "/" {
		t.Fatalf("expected path /, got %q", after.Path)
	}
	if before.Version == "" || after.Version == "" {
		t.Fatalf("expected non-empty versions before=%q after=%q", before.Version, after.Version)
	}
	if before.Version == after.Version {
		t.Fatalf("expected changed version after file write, got %q", after.Version)
	}
}

func TestFilesVersionStableWhenDirectoryContentsDoNotChange(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "note.txt", "hello")
	h := testHandlersWithRoot(t, root)
	app := fiber.New()
	app.Get("/partials/files/version", h.FilesVersion)

	before := filesVersionFor(t, app, "/")
	after := filesVersionFor(t, app, "/")

	if before.Version != after.Version {
		t.Fatalf("expected stable version, got before=%q after=%q", before.Version, after.Version)
	}
}

func TestFilesVersionChangesWhenFileContentChanges(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "note.txt", "hello")
	h := testHandlersWithRoot(t, root)
	app := fiber.New()
	app.Get("/partials/files/version", h.FilesVersion)

	before := filesVersionFor(t, app, "/")
	writeTestFile(t, root, "note.txt", "hello again")
	after := filesVersionFor(t, app, "/")

	if before.Version == after.Version {
		t.Fatalf("expected changed version after file content update, got %q", after.Version)
	}
}

func TestFilesVersionChangesWhenFileRenames(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "note.txt", "hello")
	h := testHandlersWithRoot(t, root)
	app := fiber.New()
	app.Get("/partials/files/version", h.FilesVersion)

	before := filesVersionFor(t, app, "/")
	if err := os.Rename(
		filepath.Join(root, "note.txt"),
		filepath.Join(root, "renamed.txt"),
	); err != nil {
		t.Fatalf("rename file: %v", err)
	}
	after := filesVersionFor(t, app, "/")

	if before.Version == after.Version {
		t.Fatalf("expected changed version after rename, got %q", after.Version)
	}
}

func TestFilesVersionChangesWhenFileDeletes(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "note.txt", "hello")
	h := testHandlersWithRoot(t, root)
	app := fiber.New()
	app.Get("/partials/files/version", h.FilesVersion)

	before := filesVersionFor(t, app, "/")
	if err := os.Remove(filepath.Join(root, "note.txt")); err != nil {
		t.Fatalf("remove file: %v", err)
	}
	after := filesVersionFor(t, app, "/")

	if before.Version == after.Version {
		t.Fatalf("expected changed version after delete, got %q", after.Version)
	}
}

func TestFilesRendersTrashEntryAtRootInWriteMode(t *testing.T) {
	root := t.TempDir()
	h := writeModeHandlers(t, root)

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
	if !strings.Contains(body, `data-entry-path="/trash"`) {
		t.Fatalf("expected trash entry in root listing, got %q", body)
	}
	if !strings.Contains(body, `>Trash/</span>`) {
		t.Fatalf("expected trash label in root listing, got %q", body)
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
		`data-entry-id="/note.txt"`,
		`data-entry-name="note.txt"`,
		`data-entry-path="/note.txt"`,
		`data-entry-select`,
		`data-entry-menu-button`,
		`data-panel-normal-actions`,
		`data-panel-bulk-actions`,
		`data-bulk-action="download"`,
		`class="entry-menu-heading" scope="col" aria-label="Actions"></th>`,
		`class="ui-icon"`,
		`<circle cx="12" cy="12" r="1.5"></circle>`,
		`data-relative-path="note.txt"`,
		`data-absolute-path=`,
		`data-write-context-menu`,
		`data-has-entry="false"`,
		`data-bulk="false"`,
		`data-bulk-count="0"`,
		`data-bulk-required`,
		`data-entry-required`,
		`data-single-entry-required`,
		`data-context-copy="relative"`,
		`<time datetime=`,
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
		`<th scope="col">Actions</th>`,
		`data-copy-text=`,
		`data-write-create=`,
		`data-context-write=`,
		`data-hidden-toggle`,
		`name="hidden"`,
		`>Hidden files<`,
	} {
		if strings.Contains(body, removed) {
			t.Fatalf("expected %q to be removed from table, got %q", removed, body)
		}
	}
}
