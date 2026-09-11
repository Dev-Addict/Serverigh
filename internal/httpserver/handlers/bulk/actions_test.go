package bulk

import (
	"net/http"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestBulkWriteActionMovesFilesToTrash(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "a.txt", "a")
	writeTestFile(t, root, "b.txt", "b")
	h := writeModeHandlers(t, root)
	app := fiber.New()
	app.Post("/actions/bulk/delete", h.BulkDelete)

	resp, body := formRequest(t, app, "/actions/bulk/delete", url.Values{
		"path":   {"/"},
		"target": {"/a.txt", "/b.txt"},
	})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.StatusCode, body)
	}
	assertMissing(t, filepath.Join(root, "a.txt"))
	assertMissing(t, filepath.Join(root, "b.txt"))
}

func TestBulkWriteActionDuplicatesFiles(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "a.txt", "a")
	writeTestFile(t, root, "b.txt", "b")
	h := writeModeHandlers(t, root)
	app := fiber.New()
	app.Post("/actions/bulk/duplicate", h.BulkDuplicate)

	resp, body := formRequest(t, app, "/actions/bulk/duplicate", url.Values{
		"path":   {"/"},
		"target": {"/a.txt", "/b.txt"},
	})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.StatusCode, body)
	}
	assertFileContent(t, filepath.Join(root, "a 1.txt"), "a")
	assertFileContent(t, filepath.Join(root, "b 1.txt"), "b")
}

func TestBulkWriteActionRejectsEmptySelection(t *testing.T) {
	root := t.TempDir()
	h := writeModeHandlers(t, root)
	app := fiber.New()
	app.Post("/actions/bulk/delete", h.BulkDelete)

	resp, _ := formRequest(t, app, "/actions/bulk/delete", url.Values{
		"path": {"/"},
	})

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestBulkWriteActionValidatesBeforeMutating(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "a.txt", "a")
	h := writeModeHandlers(t, root)
	app := fiber.New()
	app.Post("/actions/bulk/delete", h.BulkDelete)

	resp, _ := formRequest(t, app, "/actions/bulk/delete", url.Values{
		"path":   {"/"},
		"target": {"/a.txt", "/missing.txt"},
	})

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
	assertFileContent(t, filepath.Join(root, "a.txt"), "a")
}
