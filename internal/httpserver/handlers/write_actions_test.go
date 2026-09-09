package handlers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/config"
)

func TestWriteActionRejectsReadOnlyMode(t *testing.T) {
	h := testHandlers(t)
	app := fiber.New()
	app.Post("/actions/mkdir", h.Mkdir)

	resp, body := formRequest(t, app, "/actions/mkdir", url.Values{
		"path": {"/"},
		"name": {"docs"},
	})

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", resp.StatusCode, body)
	}
	if !strings.Contains(body, `data-error-code="write_disabled"`) {
		t.Fatalf("expected write disabled error, got %q", body)
	}
}

func TestWriteActionCreatesDirectoryAndRefreshesListing(t *testing.T) {
	root := t.TempDir()
	h := writeModeHandlers(t, root)
	app := fiber.New()
	app.Post("/actions/mkdir", h.Mkdir)

	resp, body := formRequest(t, app, "/actions/mkdir", url.Values{
		"path": {"/"},
		"name": {"docs"},
	})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.StatusCode, body)
	}
	if !strings.Contains(body, "docs/") {
		t.Fatalf("expected refreshed listing, got %q", body)
	}
	if _, err := os.Stat(filepath.Join(root, "docs")); err != nil {
		t.Fatalf("expected directory to be created: %v", err)
	}
}

func TestWriteActionCreatesFileAndRefreshesListing(t *testing.T) {
	root := t.TempDir()
	h := writeModeHandlers(t, root)
	app := fiber.New()
	app.Post("/actions/file", h.CreateFile)

	resp, body := formRequest(t, app, "/actions/file", url.Values{
		"path": {"/"},
		"name": {"note.txt"},
	})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.StatusCode, body)
	}
	if !strings.Contains(body, "note.txt") {
		t.Fatalf("expected refreshed listing, got %q", body)
	}
	content, err := os.ReadFile(filepath.Join(root, "note.txt"))
	if err != nil {
		t.Fatalf("read created file: %v", err)
	}
	if string(content) != "" {
		t.Fatalf("expected empty created file, got %q", content)
	}
}

func TestWriteActionMovesFileToTrash(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "note.txt", "hello")
	h := writeModeHandlers(t, root)
	app := fiber.New()
	app.Post("/actions/delete", h.Delete)

	resp, body := formRequest(t, app, "/actions/delete", url.Values{
		"path":   {"/"},
		"target": {"/note.txt"},
	})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.StatusCode, body)
	}
	if _, err := os.Stat(filepath.Join(root, "note.txt")); !os.IsNotExist(err) {
		t.Fatalf("expected original file to be moved, got %v", err)
	}
	trashEntries, err := os.ReadDir(filepath.Join(root, ".serverigh", "trash"))
	if err != nil {
		t.Fatalf("read trash: %v", err)
	}
	if len(trashEntries) != 1 {
		t.Fatalf("expected one trash entry, got %d", len(trashEntries))
	}
}

func TestWriteActionDuplicatesFile(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "note.txt", "hello")
	h := writeModeHandlers(t, root)
	app := fiber.New()
	app.Post("/actions/duplicate", h.Duplicate)

	resp, body := formRequest(t, app, "/actions/duplicate", url.Values{
		"path":   {"/"},
		"target": {"/note.txt"},
	})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.StatusCode, body)
	}
	content, err := os.ReadFile(filepath.Join(root, "note 1.txt"))
	if err != nil {
		t.Fatalf("read duplicate file: %v", err)
	}
	if string(content) != "hello" {
		t.Fatalf("expected duplicate content, got %q", content)
	}
}

func TestWriteActionUploadsFile(t *testing.T) {
	root := t.TempDir()
	h := writeModeHandlers(t, root)
	app := fiber.New()
	app.Post("/actions/upload", h.Upload)

	resp, body := uploadRequest(t, app, "/actions/upload", "note.txt", "hello")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.StatusCode, body)
	}
	content, err := os.ReadFile(filepath.Join(root, "note.txt"))
	if err != nil {
		t.Fatalf("read uploaded file: %v", err)
	}
	if string(content) != "hello" {
		t.Fatalf("expected uploaded content, got %q", content)
	}
}

func TestWriteActionUploadsFolderFiles(t *testing.T) {
	root := t.TempDir()
	h := writeModeHandlers(t, root)
	app := fiber.New()
	app.Post("/actions/upload", h.Upload)

	resp, body := uploadRequest(
		t,
		app,
		"/actions/upload",
		"folder/note.txt",
		"hello",
	)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.StatusCode, body)
	}
	content, err := os.ReadFile(filepath.Join(root, "folder", "note.txt"))
	if err != nil {
		t.Fatalf("read uploaded folder file: %v", err)
	}
	if string(content) != "hello" {
		t.Fatalf("expected uploaded content, got %q", content)
	}
}

func writeModeHandlers(t *testing.T, root string) Handlers {
	t.Helper()

	return testHandlersWithRoot(t, root, func(cfg *config.Config) {
		cfg.Write = true
	})
}

func formRequest(
	t *testing.T,
	app *fiber.App,
	target string,
	values url.Values,
) (*http.Response, string) {
	t.Helper()

	req := httptest.NewRequest(
		http.MethodPost,
		target,
		strings.NewReader(values.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")

	return testFiberRequest(t, app, req)
}

func uploadRequest(
	t *testing.T,
	app *fiber.App,
	target string,
	name string,
	content string,
) (*http.Response, string) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := writer.WriteField("path", "/"); err != nil {
		t.Fatalf("write path field: %v", err)
	}
	if err := writer.WriteField("relative_path", name); err != nil {
		t.Fatalf("write relative path field: %v", err)
	}
	file, err := writer.CreateFormFile("file", filepath.Base(name))
	if err != nil {
		t.Fatalf("create upload field: %v", err)
	}
	if _, err := file.Write([]byte(content)); err != nil {
		t.Fatalf("write upload content: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, target, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("HX-Request", "true")

	return testFiberRequest(t, app, req)
}

func testFiberRequest(
	t *testing.T,
	app *fiber.App,
	req *http.Request,
) (*http.Response, string) {
	t.Helper()

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("test request: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}

	return resp, string(body)
}
