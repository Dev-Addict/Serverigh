package bulk

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestBulkWriteActionRejectsDuplicateTargets(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "a.txt", "a")
	h := writeModeHandlers(t, root)
	app := fiber.New()
	app.Post("/actions/bulk/delete", h.BulkDelete)

	resp, _ := formRequest(t, app, "/actions/bulk/delete", url.Values{
		"path":   {"/"},
		"target": {"/a.txt", "/a.txt"},
	})

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
	assertFileContent(t, filepath.Join(root, "a.txt"), "a")
}

func TestBulkMoveRejectsDuplicateTargetNames(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "one/a.txt", "a")
	writeTestFile(t, root, "two/a.txt", "b")
	if err := os.Mkdir(filepath.Join(root, "dest"), 0o755); err != nil {
		t.Fatalf("create destination: %v", err)
	}
	h := writeModeHandlers(t, root)
	app := fiber.New()
	app.Post("/actions/bulk/move", h.BulkMove)

	resp, _ := formRequest(t, app, "/actions/bulk/move", url.Values{
		"destination": {"/dest"},
		"path":        {"/"},
		"target":      {"/one/a.txt", "/two/a.txt"},
	})

	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", resp.StatusCode)
	}
	assertFileContent(t, filepath.Join(root, "one/a.txt"), "a")
	assertFileContent(t, filepath.Join(root, "two/a.txt"), "b")
}
