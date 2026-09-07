package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/config"
)

func TestPreviewRendersBoundedFilePreview(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "note.txt", "abcdef")
	h := testHandlersWithRoot(t, root, func(cfg *config.Config) {
		cfg.MaxPreviewBytes = 3
	})

	app := fiber.New()
	app.Get("/preview", h.Preview)

	resp, body := testRequest(t, app, http.MethodGet, "/preview?path=/note.txt")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, "abc") {
		t.Fatalf("expected preview content, got %q", body)
	}

	if strings.Contains(body, "abcdef") {
		t.Fatalf("expected bounded preview, got %q", body)
	}

	if !strings.Contains(body, "Preview truncated") {
		t.Fatalf("expected truncation notice, got %q", body)
	}

	if !strings.Contains(body, `class="preview-chrome"`) {
		t.Fatalf("expected fixed preview chrome, got %q", body)
	}

	if !strings.Contains(body, `class="preview-body preview-body-text"`) {
		t.Fatalf("expected scrollable preview body, got %q", body)
	}
}

func TestPreviewRejectsDirectory(t *testing.T) {
	h := testHandlers(t)
	app := fiber.New()
	app.Get("/preview", h.Preview)

	resp, _ := testRequest(t, app, http.MethodGet, "/preview?path=/")

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestPreviewReturnsOperationalErrorCode(t *testing.T) {
	h := testHandlers(t)
	app := fiber.New()
	app.Get("/preview", h.Preview)

	resp, body := testRequest(t, app, http.MethodGet, "/preview?path=/")

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, `"code":"is_directory"`) {
		t.Fatalf("expected is directory code, got %q", body)
	}
}

func TestPreviewEscapesActionLinks(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "a?b#c&d.txt", "hello")
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/preview", h.Preview)

	resp, body := testRequest(
		t,
		app,
		http.MethodGet,
		"/preview?path=%2Fa%3Fb%23c%26d.txt",
	)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	normalizedBody := strings.ToLower(body)
	if !strings.Contains(normalizedBody, "/raw?path=%2fa%3fb%23c%26d.txt") {
		t.Fatalf("expected escaped raw path, got %q", body)
	}

	if !strings.Contains(
		normalizedBody,
		"/download?path=%2fa%3fb%23c%26d.txt",
	) {
		t.Fatalf("expected escaped download path, got %q", body)
	}
}

func TestPreviewRendersMarkdownHTML(t *testing.T) {
	root := t.TempDir()
	writeTestFile(
		t,
		root,
		"README.md",
		"# Title\n\n| Name | Value |\n| --- | --- |\n| One | Two |\n\n<script>x</script>",
	)
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/preview", h.Preview)

	resp, body := testRequest(t, app, http.MethodGet, "/preview?path=/README.md")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, `<article class="markdown-preview">`) {
		t.Fatalf("expected markdown preview, got %q", body)
	}

	if !strings.Contains(body, "<table>") {
		t.Fatalf("expected markdown table, got %q", body)
	}

	if strings.Contains(body, "<script>x</script>") {
		t.Fatalf("expected raw html to be skipped, got %q", body)
	}
}

func TestPreviewRendersCSVTable(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "data.csv", "name,value\none,\"two, too\"\n")
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/preview", h.Preview)

	resp, body := testRequest(t, app, http.MethodGet, "/preview?path=/data.csv")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, `class="csv-preview"`) {
		t.Fatalf("expected csv preview, got %q", body)
	}

	if !strings.Contains(body, "<td>two, too</td>") {
		t.Fatalf("expected parsed csv cell, got %q", body)
	}
}

func TestPreviewRendersPrettyJSON(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "data.json", `{"name":"one","count":2}`)
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/preview", h.Preview)

	resp, body := testRequest(t, app, http.MethodGet, "/preview?path=/data.json")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, "&#34;name&#34;") {
		t.Fatalf("expected escaped pretty json, got %q", body)
	}

	if !strings.Contains(body, `class="syntax-preview code-preview"`) {
		t.Fatalf("expected highlighted json preview, got %q", body)
	}
}

func TestPreviewRendersHighlightedCode(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "main.go", "package main\n\nfunc main() {}\n")
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/preview", h.Preview)

	resp, body := testRequest(t, app, http.MethodGet, "/preview?path=/main.go")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, `class="preview-body preview-body-code"`) {
		t.Fatalf("expected code preview body, got %q", body)
	}

	if !strings.Contains(body, `class="syntax-preview code-preview"`) {
		t.Fatalf("expected highlighted code wrapper, got %q", body)
	}

	if !strings.Contains(body, "<span") {
		t.Fatalf("expected highlighted spans, got %q", body)
	}
}

func TestPreviewRendersImagePreview(t *testing.T) {
	root := t.TempDir()
	writeBytesTestFile(
		t,
		root,
		"image.png",
		[]byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR"),
	)
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/preview", h.Preview)

	resp, body := testRequest(t, app, http.MethodGet, "/preview?path=/image.png")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, `class="media-source image-source"`) {
		t.Fatalf("expected image media preview, got %q", body)
	}

	if !strings.Contains(body, `src="/raw?path=%2Fimage.png"`) {
		t.Fatalf("expected raw image source, got %q", body)
	}
}

func TestPreviewHidesReadLimitNoticeForLargeImage(t *testing.T) {
	root := t.TempDir()
	writeBytesTestFile(
		t,
		root,
		"image.png",
		[]byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDRlarge image body"),
	)
	h := testHandlersWithRoot(t, root, func(cfg *config.Config) {
		cfg.MaxPreviewBytes = 3
	})

	app := fiber.New()
	app.Get("/preview", h.Preview)

	resp, body := testRequest(t, app, http.MethodGet, "/preview?path=/image.png")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if strings.Contains(body, "Preview truncated") {
		t.Fatalf("expected no bounded read notice for image preview, got %q", body)
	}

	if !strings.Contains(body, `class="preview-body preview-body-image"`) {
		t.Fatalf("expected image preview body, got %q", body)
	}
}

func TestPreviewHidesReadLimitNoticeForLargeVideo(t *testing.T) {
	root := t.TempDir()
	writeBytesTestFile(t, root, "movie.mp4", []byte{0x00, 0x01, 0x02, 0x03})
	h := testHandlersWithRoot(t, root, func(cfg *config.Config) {
		cfg.MaxPreviewBytes = 3
	})

	app := fiber.New()
	app.Get("/preview", h.Preview)

	resp, body := testRequest(t, app, http.MethodGet, "/preview?path=/movie.mp4")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if strings.Contains(body, "Preview truncated") {
		t.Fatalf("expected no bounded read notice for video preview, got %q", body)
	}

	if !strings.Contains(body, `class="preview-body preview-body-video"`) {
		t.Fatalf("expected video preview body, got %q", body)
	}
}

func TestPreviewRendersBinaryDetails(t *testing.T) {
	root := t.TempDir()
	writeBytesTestFile(t, root, "archive.bin", []byte{0x00, 0x01, 0x02})
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/preview", h.Preview)

	resp, body := testRequest(t, app, http.MethodGet, "/preview?path=/archive.bin")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(body, `class="binary-preview"`) {
		t.Fatalf("expected binary preview, got %q", body)
	}

	if !strings.Contains(body, "Preview is not available") {
		t.Fatalf("expected binary message, got %q", body)
	}
}
