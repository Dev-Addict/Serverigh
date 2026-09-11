package bulk

import (
	"archive/zip"
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestBulkDownloadReturnsZip(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "a.txt", "a")
	writeTestFile(t, root, "docs/b.txt", "b")
	h := testHandlersWithRoot(t, root)
	app := fiber.New()
	app.Post("/download/bulk", h.BulkDownload)

	resp, body := formRequest(t, app, "/download/bulk", url.Values{
		"target": {"/a.txt", "/docs"},
	})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.StatusCode, body)
	}
	if resp.Header.Get(fiber.HeaderContentType) != "application/zip" {
		t.Fatalf("expected zip content type")
	}

	reader, err := zip.NewReader(
		bytes.NewReader([]byte(body)),
		int64(len(body)),
	)
	if err != nil {
		t.Fatalf("read zip: %v", err)
	}
	assertZipEntries(t, reader, "a.txt", "docs/b.txt")
}

func TestBulkDownloadRejectsEmptySelection(t *testing.T) {
	h := testHandlers(t)
	app := fiber.New()
	app.Post("/download/bulk", h.BulkDownload)

	resp, body := formRequest(t, app, "/download/bulk", url.Values{})

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", resp.StatusCode, body)
	}
	if resp.Header.Get(fiber.HeaderContentType) == "application/zip" {
		t.Fatalf("expected validation before zip headers")
	}
}

func TestBulkDownloadValidatesBeforeStreaming(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "a.txt", "a")
	h := testHandlersWithRoot(t, root)
	app := fiber.New()
	app.Post("/download/bulk", h.BulkDownload)

	resp, _ := formRequest(t, app, "/download/bulk", url.Values{
		"target": {"/a.txt", "/missing.txt"},
	})

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
	if resp.Header.Get(fiber.HeaderContentType) == "application/zip" {
		t.Fatalf("expected validation before zip headers")
	}
}
